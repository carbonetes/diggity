package curator

import (
	"fmt"

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

	// tag, err := name.NewTag("ubuntu:new")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	image, err := tarball.ImageFromPath(path, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(image)
	return image, nil
}
