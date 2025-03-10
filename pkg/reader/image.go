package reader

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/carbonetes/diggity/cmd/diggity/config"
	stream "github.com/carbonetes/diggity/cmd/diggity/grove"
	"github.com/carbonetes/diggity/cmd/diggity/ui"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/scanner"
	"github.com/carbonetes/diggity/pkg/types"
	"github.com/golistic/urn"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/uuid"
)

// GetImage retrieves a Docker image given its name or digest.
// It first checks if the image exists locally, and if not, it pulls it from the remote registry.
// Returns the image object and an error if any.
func GetImage(input string, config *types.RegistryConfig) (*v1.Image, *name.Reference, error) {
	ref, err := name.ParseReference(input)
	if err != nil {
		return nil, nil, err
	}

	var image v1.Image
	exists, image, _ := CheckIfImageExistsInLocal(ref)
	if exists {
		return &image, &ref, nil
	}

	if config != nil {
		// Load image remotely with authn.Basic
		// Check out information about authn in https://github.com/google/go-containerregistry/tree/main/pkg/authn
		image, err = remote.Image(ref, remote.WithAuth(&authn.Basic{
			Username: config.Username,
			Password: config.Password,
		}))
		if err != nil {
			if strings.Contains(err.Error(), "unsupported MediaType") {
				return nil, nil, errors.New(fmt.Sprint(ErrUnsupportedMediaType))

			} else if strings.Contains(err.Error(), "authentication required") {
				return nil, nil, errors.New(ErrNotExistOrAuthenticationRequired)
			} else {
				return nil, nil, err
			}
		}

	} else {
		// Remotely load image from public registry
		image, err = remote.Image(ref)
		if err != nil {
			return nil, nil, err
		}
	}

	return &image, &ref, nil
}

// ReadFiles reads the layers of a given v1.Image and processes its contents.
// It returns an error if there's any issue encountered while reading the layers or processing its contents.
func ReadFiles(image *v1.Image, addr *urn.URN) error {
	if image == nil {
		return errors.New("image is nil")
	}

	layers, err := (*image).Layers()
	if err != nil {
		return err
	}

	// Get the maximum file size from the configuration file
	// If maxFileSize is not set, use the default value of 50MB
	maxFileSize := config.Config.MaxFileSize

	var wg sync.WaitGroup
	wg.Add(len(layers))
	for _, layer := range layers {
		hash, err := layer.Digest()
		if err != nil {
			log.Debugf("failed to get layer digest: %s", err)
		}
		go func(layer v1.Layer) {
			defer wg.Done()
			contents, err := layer.Uncompressed()
			if err != nil {
				log.Debugf("Failed to uncompress layer: %s", err)
			}

			err = processLayerContents(hash.String(), contents, maxFileSize, addr)
			if err != nil {
				log.Debugf("Failed to process layer contents: %s", err)
			}
		}(layer)

	}
	wg.Wait()
	return nil
}

// processLayerContents reads the contents of a tar file and processes each file header.
// It skips files that exceed the maximum file size and only processes regular files.
// The processed files are hashed using the layerHash and stored for later use.
// Returns an error if there was an issue reading or processing the tar file.
func processLayerContents(layer string, contents io.ReadCloser, maxFileSize int64, addr *urn.URN) error {
	defer contents.Close()
	reader := tar.NewReader(contents)
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		ui.AddFile(header.Name)

		if err := processHeader(header, reader, layer, maxFileSize, addr); err != nil {
			log.Debug(err)
		}
	}

	return nil
}

func processHeader(header *tar.Header, reader io.Reader, layer string, maxFileSize int64, addr *urn.URN) error {
	if slices.Contains(archiveTypes, filepath.Ext(header.Name)) {
		return processArchiveFileFromLayer(header, reader, addr)
	}

	if header.Typeflag == tar.TypeReg {
		if header.Size > maxFileSize {
			return nil
		}
		return processTarHeader(layer, header, reader, addr)
	}

	return nil
}

func processArchiveFileFromLayer(header *tar.Header, reader io.Reader, addr *urn.URN) error {
	b, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	processArchive(bytes.NewReader(b), header.Name, header.Size, addr)
	return nil
}

// processTarHeader processes a tar header and its contents, checking if the file size is within the limit and if it is a regular file.
// If the file is a related file, it processes the file and returns an error if encountered.
func processTarHeader(layer string, header *tar.Header, reader io.Reader, addr *urn.URN) error {
	category, matched, readFlag := scanner.CheckRelatedFiles(header.Name)
	if matched {
		if !readFlag {
			stream.Emit(category, types.Payload{
				Address: addr,
				Layer:   layer,
				Body:    header.Name,
			})
			return nil
		}
		err := processFile(header.Name, layer, reader, category, addr)
		if err != nil {
			return err
		}
	}

	return nil
}

// processFile reads the contents of a file from a reader, creates a temporary file,
// writes the contents of the reader to the temporary file, reads the content of the
// temporary file and emits a manifest file to a stream.
// The manifest file contains the name of the file, the hash of the layer and the content of the file.
// It returns an error if any of the operations fail.
func processFile(name, layer string, reader io.Reader, category string, addr *urn.URN) error {
	f, err := os.Create(os.TempDir() + string(os.PathSeparator) + "diggity-tmp-" + uuid.NewString())
	if err != nil {
		return err
	}

	_, err = io.Copy(f, reader)
	if err != nil {
		return err
	}

	if category == "rpm" {
		err = handleRpmFile(f.Name(), category, layer, addr)
		if err != nil {
			return err
		}
	} else {
		err = handleManifestFile(name, category, layer, f, addr)
		if err != nil {
			return err
		}
	}

	f.Close()
	os.Remove(f.Name())

	return nil
}

// CheckIfImageExistsInLocal checks if the given image reference exists locally.
// It returns a boolean indicating whether the image exists, the image object if it exists, and an error if any.
func CheckIfImageExistsInLocal(ref name.Reference) (bool, v1.Image, error) {
	img, err := daemon.Image(ref)
	if err != nil {
		// If the error indicates the image is not found, return false.
		if strings.Contains(err.Error(), "not found") {
			return false, nil, nil
		}
		// For other errors, return the error.
		return false, nil, err
	}

	// If no error, the image is found locally.
	return true, img, nil
}
