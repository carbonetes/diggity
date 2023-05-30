package dart

import (
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/google/uuid"
)

func newPackage(dartMetadata interface{}) *model.Package {
	var pkg model.Package
	var licenses []string
	pkg.ID = uuid.NewString()
	switch md := dartMetadata.(type) {
	case Metadata:
		pkg.Name = md["name"].(string)
		pkg.Type = Type
		pkg.Path = md["name"].(string)
		//check if version exist, if not set default of 0.0.0
		if val, ok := md["version"].(string); ok {
			pkg.Version = val
		} else {
			pkg.Version = "0.0.0"
		}

		if val, ok := md["description"].(string); ok {
			pkg.Description = val
		}

		if val, ok := md["license"].(string); ok {
			licenses = append(licenses, val)
		}
		pkg.Licenses = licenses
		pkg.Metadata = md
		generateCpes(&pkg, &md)
	case metadata.PubspecLockMetadata:
		pkg.Name = md.Description.Name
		pkg.Version = md.Version
		pkg.Type = Type
		pkg.Path = md.Description.Name
		pkg.Metadata = md
		generateCpes(&pkg, nil)
	}
	setPurl(&pkg)

	return &pkg
}

// Parse PURL
func setPurl(pkg *model.Package) {
	pkg.PURL = model.PURL("pkg" + ":" + "dart" + "/" + pkg.Name + "@" + pkg.Version)
}
