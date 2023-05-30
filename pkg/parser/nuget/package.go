package nuget

import (
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/google/uuid"
)

func newPackage(keyValue string, lib *metadata.DotnetLibrary) *model.Package {
	var pkg model.Package
	properties := strings.Split(keyValue, "/")
	if len(properties) != 2 {
		return nil
	}

	name, version := properties[0], properties[1]
	pkg.ID = uuid.NewString()
	pkg.Name = name
	pkg.Version = version
	pkg.Type = Type
	pkg.Path = name

	setPurl(&pkg)
	generateCpes(&pkg)
	pkg.Metadata = lib
	return &pkg
}

// Parse PURL
func setPurl(pkg *model.Package) {
	pkg.PURL = model.PURL("pkg" + ":" + "dotnet" + "/" + pkg.Name + "@" + pkg.Version)
}
