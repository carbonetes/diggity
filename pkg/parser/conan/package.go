package conan

import (
	"log"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/util"
	"github.com/google/uuid"
)

func newPackage(location *model.Location, conanMetadata interface{}) *model.Package {
	pkg := new(model.Package)
	pkg.ID = uuid.NewString()

	var name, version string
	switch md := conanMetadata.(type) {
	case metadata.ConanMetadata:
		name, version = md.Name, md.Version
		pkg.Metadata = md
		log.Print(name)
	case metadata.ConanLockNode:
		name, version = parseRef(md.Ref)
		pkg.Metadata = md
	}

	pkg.Name = name
	pkg.Version = version
	pkg.Path = pkg.Name
	pkg.Type = Type
	pkg.Locations = append(pkg.Locations, model.Location{
		Path:      util.TrimUntilLayer(*location),
		LayerHash: location.LayerHash,
	})
	pkg.Licenses = []string{}
	setPurl(pkg)
	generateCpes(pkg)

	return pkg
}

func setPurl(pkg *model.Package) {
	pkg.PURL = model.PURL("pkg" + ":" + Type + "/" + pkg.Name + "@" + pkg.Version)
}

func parseRef(ref string) (name string, version string) {
	var nv []string
	if strings.Contains(ref, "@") {
		nv = strings.Split(ref, "@")
	} else if strings.Contains(ref, "#") {
		nv = strings.Split(ref, "#")
	} else {
		nv = append(nv, ref)
	}

	result := strings.Split(nv[0], "/")
	name = result[0]
	version = result[1]

	if strings.Contains(version, "[") && strings.Contains(version, "]") {
		version = strings.Replace(version, "[", "", -1)
		version = strings.Replace(version, "]", "", -1)
	}

	return name, version
}
