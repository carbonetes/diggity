package docker

import (
	"fmt"
	"time"

	"github.com/carbonetes/diggity/cmd/diggity/config"
)

// TarballInfo contains information about a Docker tarball
type TarballInfo struct {
	Path         string    `json:"path"`
	MediaType    string    `json:"media_type"`
	LayerCount   int       `json:"layer_count"`
	TotalSize    int64     `json:"total_size"`
	Created      time.Time `json:"created"`
	Architecture string    `json:"architecture"`
	OS           string    `json:"os"`
}

// LayerInfo contains information about a Docker layer
type LayerInfo struct {
	Digest    string `json:"digest"`
	Size      int64  `json:"size"`
	MediaType string `json:"media_type"`
}

// ImageAnalysis contains comprehensive analysis of a Docker image
type ImageAnalysis struct {
	MediaType    string      `json:"media_type"`
	Architecture string      `json:"architecture"`
	OS           string      `json:"os"`
	Created      time.Time   `json:"created"`
	LayerCount   int         `json:"layer_count"`
	TotalSize    int64       `json:"total_size"`
	Layers       []LayerInfo `json:"layers"`
	Env          []string    `json:"env,omitempty"`
	Cmd          []string    `json:"cmd,omitempty"`
	Entrypoint   []string    `json:"entrypoint,omitempty"`
}

var (
	// ErrUnsupportedMediaType is the error message for unsupported media type
	ErrUnsupportedMediaType = "Error: Unsupported MediaType Detected\n\nThis issue is often encountered when interacting with older image manifests or registries that have not been updated to support the current Docker distribution specifications. Please consider upgrading your container registry or converting your image manifests to a supported version. For more information and potential workarounds, refer to the discussion at https://github.com/google/go-containerregistry/issues/377.\n"

	// ErrNotExistOrAuthenticationRequired is the error message for authentication required
	ErrNotExistOrAuthenticationRequired = fmt.Sprintf("Error: Image Not Found or Authentication Required\n\nThe target image may not exist, or the registry requires authentication to access the image. Please provide the required credentials to authenticate with the registry by editing %s.", config.GetConfigPath())

	// ErrImageNotFound indicates that the Docker image was not found
	ErrImageNotFound = fmt.Errorf("docker image not found")

	// ErrTarballInvalid indicates that the tarball file is invalid or corrupted
	ErrTarballInvalid = fmt.Errorf("invalid or corrupted docker tarball")

	// ErrLayerProcessing indicates an error occurred while processing Docker layers
	ErrLayerProcessing = fmt.Errorf("error processing docker layer")

	// ErrRegistryConnection indicates a connection error to the Docker registry
	ErrRegistryConnection = fmt.Errorf("failed to connect to docker registry")
)
