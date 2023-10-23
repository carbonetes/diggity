package pnpm

import (
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/language/javascript/npm"
	"github.com/google/uuid"
)

func newPackage(name string, version string) *model.Package {
	return &model.Package{
		ID:      uuid.NewString(),
		Name:    strings.ReplaceAll(name, "@", ""),
		PackageOrigin: model.ApplicationPackage,
		Language: npm.Language,
		Parser: Type,
		Version: version,
		Type:    Type,
		PURL:    setPurl(name, version),
	}
}

// Parse PURL
func setPurl(name string, version string) model.PURL {
	return model.PURL("pkg" + ":" + Type + "/" + name + "@" + version)
}
