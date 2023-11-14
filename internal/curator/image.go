package curator

import (
	"strings"

	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/remote"
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
	stream.SetImageInstance(image)
	return image, nil
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
