package reader

import (
	"fmt"

	"github.com/carbonetes/diggity/cmd/diggity/config"
)

var (
	// ErrUnsupportedMediaType is the error message for unsupported media type.
	ErrUnsupportedMediaType = "Error: Unsupported MediaType Detected\n\nThis issue is often encountered when interacting with older image manifests or registries that have not been updated to support the current Docker distribution specifications. Please consider upgrading your container registry or converting your image manifests to a supported version. For more information and potential workarounds, refer to the discussion at https://github.com/google/go-containerregistry/issues/377.\n"

	// ErrAuthenticationRequired is the error message for authentication required.
	ErrNotExistOrAuthenticationRequired = fmt.Sprintf("Error: Image Not Found or Authentication Required\n\nThe target image may not exist, or the registry requires authentication to access the image. Please provide the required credentials to authenticate with the registry by editing %s.", config.GetConfigPath())
)
