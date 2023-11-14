package curator

import (
	"archive/tar"
	"io"
	"os"

	"github.com/carbonetes/diggity/pkg/types"
	"github.com/google/uuid"
)

func ReadTarball(reader io.Reader) error {
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		f, err := os.Create(os.TempDir() + string(os.PathSeparator) + "diggity-tmp-" + uuid.NewString())
		if err != nil {
			return err
		}

		_, err = io.Copy(f, reader)
		if err != nil {
			return err
		}

		manifest := types.ManifestFile{
			Path: header.Name,
		}

		err = manifest.ReadContent(f)
		if err != nil {
			return err
		}

	}
	return nil
}
