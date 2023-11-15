package curator

import (
	"archive/tar"
	"io"
	"os"
	"strings"

	"github.com/carbonetes/diggity/internal/scanner"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/uuid"
)

// GetImage retrieves a Docker image given its name or digest.
// It first checks if the image exists locally, and if not, it pulls it from the remote registry.
// Returns the image object and an error if any.
func GetImage(input string) (v1.Image, error) {
	ref, err := name.ParseReference(input)
	if err != nil {
		return nil, err
	}
	var image v1.Image
	exists, image, err := CheckIfImageExistsInLocal(ref)
	if !exists || err != nil {
		image, err = remote.Image(ref)
		if err != nil {
			return nil, err
		}
	}

	return image, nil
}

// ReadFiles reads the layers of a given v1.Image and processes its contents.
// It returns an error if there's any issue encountered while reading the layers or processing its contents.
func ReadFiles(image v1.Image) error {
	layers, err := image.Layers()
	if err != nil {
		return err
	}
	maxFileSize := stream.GetParameters().MaxFileSize
	for _, layer := range layers {
		contents, err := layer.Uncompressed()
		if err != nil {
			return err
		}

		layerHash, err := layer.Digest()
		if err != nil {
			return err
		}

		err = processLayerContents(contents, layerHash.String(), maxFileSize)
		if err != nil {
			return err
		}
	}

	return nil
}

// processLayerContents reads the contents of a tar file and processes each file header.
// It skips files that exceed the maximum file size and only processes regular files.
// The processed files are hashed using the layerHash and stored for later use.
// Returns an error if there was an issue reading or processing the tar file.
func processLayerContents(contents io.ReadCloser, layerHash string, maxFileSize int64) error {
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
		if header.Size > maxFileSize {
			continue
		}
		if header.Typeflag == tar.TypeReg {
			err = processTarHeader(header, reader, layerHash, maxFileSize)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// processTarHeader processes a tar header and its contents, checking if the file size is within the limit and if it is a regular file.
// If the file is a related file, it processes the file and returns an error if encountered.
func processTarHeader(header *tar.Header, reader io.Reader, layerHash string, maxFileSize int64) error {
	if header.Size > maxFileSize {
		return nil
	}
	if header.Typeflag == tar.TypeReg {
		stream.Emit(stream.FilesystemCheckEvent, header.Name)
		category, matched := scanner.CheckRelatedFiles(header.Name)
		if matched {
			err := processFile(header.Name, layerHash, reader, category)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// processFile reads the contents of a file from a reader, creates a temporary file, 
// writes the contents of the reader to the temporary file, reads the content of the 
// temporary file and emits a manifest file to a stream.
// The manifest file contains the name of the file, the hash of the layer and the content of the file.
// It returns an error if any of the operations fail.
func processFile(name string, layerHash string, reader io.Reader, category string) error {
	f, err := os.Create(os.TempDir() + string(os.PathSeparator) + "diggity-tmp-" + uuid.NewString())
	if err != nil {
		return err
	}

	_, err = io.Copy(f, reader)
	if err != nil {
		return err
	}

	manifest := types.ManifestFile{
		Path:  name,
		Layer: layerHash,
	}

	err = manifest.ReadContent(f)
	if err != nil {
		return err
	}
	f.Close()
	os.Remove(f.Name())
	stream.Emit(category, manifest)

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
