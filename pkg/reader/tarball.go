package reader

import (
	"fmt"

	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

// ReadTarballAsImage reads a tarball from the given path and returns it as a v1.Image.
func ReadTarballAsImage(path string) (v1.Image, error) {
	found, err := helper.IsFileExists(path)
	if !found {
		return nil, fmt.Errorf("file not found: %s", path)
	}

	if err != nil {
		return nil, err
	}

	image, err := tarball.ImageFromPath(path, nil)
	if err != nil {
		log.Error(err)
	}

	return image, nil
}
