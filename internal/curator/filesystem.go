package curator

import (
	"archive/tar"
	"io"
	"os"

	"github.com/carbonetes/diggity/internal/scanner"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/uuid"
)

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
