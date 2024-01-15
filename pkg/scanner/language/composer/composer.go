package composer

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "composer"

var Manifests = []string{"composer.lock"}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)
	if !ok {
		log.Error("Composer Handler received unknown type")
	}

	metadata := readManifestFile(manifest.Content)

	for _, pkg := range metadata.Packages {
		if pkg.Name == "" {
			continue
		}

		component := types.NewComponent(pkg.Name, pkg.Version, Type, manifest.Path, pkg.Description, pkg)
		component.Licenses = pkg.License
		stream.AddComponent(component)
	}

	for _, pkg := range metadata.PackagesDev {
		if pkg.Name == "" {
			continue
		}

		component := types.NewComponent(pkg.Name, pkg.Version, Type, manifest.Path, pkg.Description, pkg)
		component.Licenses = pkg.License
		props := strings.Split(component.Name, "/")

		if len(props) == 0 {
			props = []string{component.Name, component.Name}
		}
		vendor, product := props[0], props[1]
		cpes := cpe.NewCPE23(vendor, product, component.Version, Type)
		if len(cpes) > 0 {
			component.CPEs = append(component.CPEs, cpes...)
		}
		stream.AddComponent(component)
	}

	return data
}
