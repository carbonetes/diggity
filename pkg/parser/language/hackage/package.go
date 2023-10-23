package hackage

import (
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/util"
	"github.com/google/uuid"
)

// Init Hackage Package
func initHackagePackage(location *model.Location, dep string, url string) *model.Package {
	name, version, pkgHash, size, rev := parseExtraDep(dep)

	pkg := new(model.Package)
	pkg.ID = uuid.NewString()
	pkg.Name = name
	pkg.Version = version
	pkg.Path = pkg.Name
	pkg.Type = Type
	pkg.PackageOrigin = model.ApplicationPackage
	pkg.Language = Language
	pkg.Parser = Type
	pkg.Locations = append(pkg.Locations, model.Location{
		Path:      util.TrimUntilLayer(*location),
		LayerHash: location.LayerHash,
	})
	pkg.Licenses = []string{}

	// get purl
	parseHackageURL(pkg)

	// get CPEs
	generateCpes(pkg)

	// fill metadata
	pkg.Metadata = metadata.HackageMetadata{
		Name:        name,
		Version:     version,
		PkgHash:     pkgHash,
		Size:        size,
		Revision:    rev,
		SnapshotURL: url,
	}

	return pkg
}

// Parse PURL
func parseHackageURL(pkg *model.Package) {
	pkg.PURL = model.PURL("pkg" + ":" + Type + "/" + pkg.Name + "@" + pkg.Version)
}

// Format Name and Version for parsing
func formatCabalPackage(anyPkg string) string {
	pkg := strings.Replace(strings.TrimSpace(anyPkg), anyTag, "", -1)
	nv := strings.Replace(pkg, " ==", "-", -1)
	return strings.Replace(nv, ",", "", -1)
}
