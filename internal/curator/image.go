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

func ReadFiles(image v1.Image) error {
	layers, err := image.Layers()
	if err != nil {
		return err
	}

	for _, layer := range layers {
		contents, err := layer.Uncompressed()
		if err != nil {
			return err
		}

		layerHash, err := layer.Digest()
		if err != nil {
			return err
		}

		err = processLayerContents(contents, layerHash.String())
		if err != nil {
			return err
		}
	}

	return nil
}

func processLayerContents(contents io.ReadCloser, layerHash string) error {
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
		if header.Typeflag == tar.TypeReg {
			category, matched := scanner.CheckRelatedFiles(header.Name)
			if matched {
				err = processFile(header.Name, layerHash, reader, category)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func processFile(name string, layerHash string, reader io.Reader, category string) error {
	f, err := os.Create(os.TempDir() + string(os.PathSeparator) + "diggity-tmp-" + uuid.NewString())
	if err != nil {
		return err
	}
	defer f.Close()

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

	stream.Emit(category, manifest)

	return nil
}

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
