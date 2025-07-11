package docker

import (
	"strings"

	"github.com/carbonetes/diggity/pkg/types"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// ImageHandler handles Docker image operations including pulling from registries
type ImageHandler struct {
	config *types.RegistryConfig
}

// NewImageHandler creates a new ImageHandler instance
func NewImageHandler(registryConfig *types.RegistryConfig) *ImageHandler {
	return &ImageHandler{
		config: registryConfig,
	}
}

// GetImage retrieves a Docker image given its name or digest
// It first checks if the image exists locally, and if not, pulls it from the remote registry
func (h *ImageHandler) GetImage(input string) (*v1.Image, *name.Reference, error) {
	ref, err := name.ParseReference(input)
	if err != nil {
		return nil, nil, err
	}

	var image v1.Image
	exists, localImage, _ := h.checkLocalImage(ref)
	if exists {
		return &localImage, &ref, nil
	}

	if h.config != nil {
		// Load image remotely with authentication
		image, err = remote.Image(ref, remote.WithAuth(&authn.Basic{
			Username: h.config.Username,
			Password: h.config.Password,
		}))
	} else {
		// Load image remotely without authentication
		image, err = remote.Image(ref)
	}

	if err != nil {
		return nil, nil, err
	}

	return &image, &ref, nil
}

// checkLocalImage checks if an image exists locally in the Docker daemon
func (h *ImageHandler) checkLocalImage(ref name.Reference) (bool, v1.Image, error) {
	image, err := daemon.Image(ref)
	if err != nil {
		if strings.Contains(err.Error(), "No such image") {
			return false, nil, nil
		}
		return false, nil, err
	}
	return true, image, nil
}

// PullImage pulls an image from a remote registry
func (h *ImageHandler) PullImage(ref name.Reference) (v1.Image, error) {
	var image v1.Image
	var err error

	if h.config != nil {
		// Load image remotely with authentication
		image, err = remote.Image(ref, remote.WithAuth(&authn.Basic{
			Username: h.config.Username,
			Password: h.config.Password,
		}))
	} else {
		// Load image remotely without authentication
		image, err = remote.Image(ref)
	}

	return image, err
}

// GetImageManifest retrieves the manifest of a Docker image
func (h *ImageHandler) GetImageManifest(image v1.Image) (*v1.Manifest, error) {
	return image.Manifest()
}

// GetImageConfig retrieves the configuration of a Docker image
func (h *ImageHandler) GetImageConfig(image v1.Image) (*v1.ConfigFile, error) {
	return image.ConfigFile()
}

// ListImageLayers returns the layers of a Docker image
func (h *ImageHandler) ListImageLayers(image v1.Image) ([]v1.Layer, error) {
	return image.Layers()
}
