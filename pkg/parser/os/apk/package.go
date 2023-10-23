package apk

import (
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/google/uuid"
)

const alpine string = "alpine"

// Parse PURL
func setPURL(pkg *model.Package) {
	arch, ok := pkg.Metadata.(Metadata)["Architecture"]
	if !ok {
		arch = ""
	}
	origin, ok := pkg.Metadata.(Metadata)["Origin"]
	if !ok {
		origin = ""
	}

	pkg.PURL = model.PURL("pkg" + `:` + Type + `/` + alpine + `/` + pkg.Name + `@` + pkg.Version + `?arch=` + arch.(string) + `&` + `upstream=` + origin.(string) + `&distro=` + alpine)
}

// Create Alpine Package
func newPackage(content string, listFiles bool) *model.Package {

	var pkg = model.Package{}

	content = strings.TrimSpace(content)
	attributes := strings.Split(content, "\n")

	metadata := parseMetadata(attributes, listFiles)
	if metadata["Name"] == nil {
		return nil
	}

	pkg.ID = uuid.NewString()
	pkg.Name = metadata["Name"].(string)
	pkg.Version = metadata["Version"].(string)
	pkg.Type = Type
	pkg.Distro = Distro
	pkg.Parser = Type
	pkg.PackageOrigin = model.OSPackage
	pkg.Description = metadata["Description"].(string)
	pkg.Metadata = metadata
	setPURL(&pkg)
	generateAlpineCpes(&pkg)

	for _, license := range strings.Split(metadata["License"].(string), " ") {
		if !strings.Contains(strings.ToLower(license), "and") {
			pkg.Licenses = append(pkg.Licenses, license)
		}
	}

	return &pkg
}
