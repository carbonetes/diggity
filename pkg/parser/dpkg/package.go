package dpkg

import (
	"path/filepath"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/google/uuid"
)

var (
	dpkgDocPath = filepath.Join("usr", "share", "doc")
	copyright   = filepath.Join("copyright")
)

func newPackage(metadata Metadata) *model.Package {
	var pkg model.Package
	if metadata["package"] == nil || metadata["version"] == nil {
		return nil
	}
	pkg.ID = uuid.NewString()
	pkg.Type = Type
	pkg.PackageOrigin = model.OSPackage
	pkg.Distro = Distro
	pkg.Parser = Type
	pkg.Name = metadata["package"].(string)
	pkg.Version = metadata["version"].(string)
	if val, ok := metadata["description"].(string); ok {
		pkg.Description = val
	}

	//need to add distro in purl
	setPurl(&pkg, metadata["architecture"].(string))

	//get CPEs
	generateCpes(&pkg)

	return &pkg
}

// Parse PURL
func setPurl(pkg *model.Package, architecture string) {
	pkg.PURL = model.PURL("pkg" + ":" + "deb" + "/" + pkg.Name + "@" + pkg.Version + "?arch=" + architecture)
}
