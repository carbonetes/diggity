package output

import "github.com/carbonetes/diggity/internal/model"

// Output final output
type Output struct {
	Packages  []*model.Package     `json:"packages"`
	Secret    *model.SecretResults `json:"secrets,omitempty"`
	ImageInfo model.ImageInfo      `json:"imageInfo"`
	Distro    *model.Distro        `json:"distro"`
}
