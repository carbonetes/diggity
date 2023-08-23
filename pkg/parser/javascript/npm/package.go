package npm

import (
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/google/uuid"
)

func newNpmPackage(metadata *metadata.PackageJSON) *model.Package {
	var pkg model.Package
	pkg.ID = uuid.NewString()

	// // init npm data
	pkg.Name = NpmMetadata.Name
	pkg.Version = NpmMetadata.Version
	pkg.Description = NpmMetadata.Description
	pkg.Type = Type
	pkg.Path = NpmMetadata.Name

	// // check type of license then parse
	switch NpmMetadata.License.(type) {
	case string:
		pkg.Licenses = append(pkg.Licenses, NpmMetadata.License.(string))
	case map[string]interface{}:
		license := NpmMetadata.License.(map[string]interface{})
		if _, ok := license["type"]; ok {
			pkg.Licenses = append(pkg.Licenses, license["type"].(string))
		}
	}

	// //parseURL
	setPurl(&pkg)
	generateCPEs(&pkg)
	pkg.Metadata = NpmMetadata

	return &pkg
}

func newNpmLockPackage(name string, metadata *metadata.LockDependency) *model.Package {
	var pkg model.Package
	pkg.ID = uuid.NewString()
	// // init npm data
	pkg.Name = name
	pkg.Version = metadata.Version
	pkg.Type = Type
	pkg.Path = name

	// //parseURL
	setPurl(&pkg)
	generateCPEs(&pkg)
	pkg.Metadata = metadata
	return &pkg
}

func newYarnLockPackage(metadata *LockMetadata) *model.Package {
	var pkg model.Package
	pkg.ID = uuid.NewString()
	pkg.Type = Type
	pkg.Name = (*metadata)["Name"].(string)
	pkg.Path = (*metadata)["Name"].(string)

	if (*metadata)["Version"] != nil {
		pkg.Version = (*metadata)["Version"].(string)
	}
	setPurl(&pkg)
	generateCPEs(&pkg)
	pkg.Metadata = metadata

	return &pkg
}

// Parse PURL
func setPurl(pkg *model.Package) {
	pkg.PURL = model.PURL("pkg" + ":" + Type + "/" + pkg.Name + "@" + pkg.Version)
}
