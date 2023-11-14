package helper

import "strings"

const (
	DefaultImageTag string = "latest"
	ImageSeparator  string = ":"
)

// FormatImage set default image tag if not provided
func FormatImage(image string) string {
	if !strings.Contains(image, ImageSeparator) {
		return image + ImageSeparator + DefaultImageTag
	}
	return image
}
