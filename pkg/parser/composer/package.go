package composer

import (
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/google/uuid"
)

func newPackage(composerPackage *metadata.ComposerPackage) *model.Package {
	var pkg model.Package
	pkg.ID = uuid.NewString()
	pkg.Name = composerPackage.Name
	pkg.Version = composerPackage.Version
	pkg.Description = composerPackage.Description
	pkg.Licenses = composerPackage.License
	pkg.Type = Type
	pkg.Path = composerPackage.Name
	setPurl(&pkg)
	generateComposerPackageCpes(&pkg)
	return &pkg
}

// Parse PURL
func setPurl(pkg *model.Package) {
	pkg.PURL = model.PURL("pkg" + ":" + "composer" + "/" + pkg.Name + "@" + pkg.Version)
}
