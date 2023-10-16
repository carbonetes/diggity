package pacman

import (
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/google/uuid"
)

func newPackage(metadata Metadata) *model.Package {
	if metadata["Name"] == nil || metadata["Name"] == "" {
		return nil
	}

	return &model.Package{
		ID:          uuid.NewString(),
		Name:        metadata["Name"].(string),
		PURL:        setPURL(metadata),
		Type:        Type,
		Version:     metadata["Version"].(string),
		Description: metadata["Desc"].(string),
		Metadata:    metadata,
	}
}

func setPURL(metadata Metadata) model.PURL {
	arch, ok := metadata["Arch"].(string)
	if !ok {
		arch = ""
	}
	return model.PURL("pkg:" + Type + "/archlinux/" + metadata["Name"].(string) + `@` + metadata["Version"].(string) + `?arch=` + arch + `&distro=` + "archlinux")
}
