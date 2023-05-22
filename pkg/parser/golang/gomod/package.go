package gomod

import (
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/golang"
	"github.com/carbonetes/diggity/pkg/parser/util"
	"github.com/google/uuid"
	"golang.org/x/mod/modfile"
)

// Initialize go mod package contents
func newPackage(location *model.Location, modPkg interface{}) *model.Package {
	var pkg model.Package
	pkg.ID = uuid.NewString()

	switch goPkg := modPkg.(type) {
	case *modfile.Require:
		pkg.Name = goPkg.Mod.Path
		pkg.Version = goPkg.Mod.Version
	case *modfile.Replace:
		pkg.Name = goPkg.New.Path
		pkg.Version = goPkg.New.Version
	}

	pkg.Type = golang.Type
	pkg.Path = pkg.Name
	pkg.Licenses = []string{}

	// get locations
	pkg.Locations = append(pkg.Locations, model.Location{
		Path:      util.TrimUntilLayer(*location),
		LayerHash: location.LayerHash,
	})

	// get CPEs
	paths := golang.SplitPath(pkg.Name)

	golang.GenerateCpes(&pkg, paths)

	// get purl
	golang.SetPurl(&pkg)

	// parse and fill final metadata
	metadata := parseMetadata(modPkg)

	if metadata != nil {
		pkg.Metadata = metadata
	}

	return &pkg
}
