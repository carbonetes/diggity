package pnpm

import (
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/google/uuid"
)

func newPackage(name string, version string) *model.Package {
	return &model.Package{
		ID:      uuid.NewString(),
		Name:    strings.ReplaceAll(name, "@", ""),
		Version: version,
		Type:    Type,
		PURL:    setPurl(name, version),
	}
}

// Parse PURL
func setPurl(name string, version string) model.PURL {
	return model.PURL("pkg" + ":" + Type + "/" + name + "@" + version)
}
