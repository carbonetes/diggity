package curator

import (
	"github.com/carbonetes/diggity/internal/helper"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

func ReadTarball(path string) (v1.Image, error) {
	found, err := helper.IsFileExists(path)
	if err != nil {
		log.Fatal(err)
	}

	if !found {
		log.Fatal("Path does not exist")
	}

	image, err := tarball.ImageFromPath(path, nil)
	if err != nil {
		log.Fatal(err)
	}
	
	return image, nil
}
