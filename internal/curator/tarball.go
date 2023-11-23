package curator

import (
	"github.com/carbonetes/diggity/internal/helper"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

// ReadTarballAsImage reads a tarball from the given path and returns it as a v1.Image.
func ReadTarballAsImage(path string) (v1.Image, error) {
	found, err := helper.IsFileExists(path)
	if err != nil {
		log.Error(err)
	}

	if !found {
		log.Fatal("Path does not exist")
	}

	image, err := tarball.ImageFromPath(path, nil)
	if err != nil {
		log.Error(err)
	}

	return image, nil
}
