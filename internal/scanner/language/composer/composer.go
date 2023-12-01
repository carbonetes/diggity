package composer

import (
	"path/filepath"
	"slices"

	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "composer"

var (
	Manifests = []string{"composer.lock"}
	log       = logger.GetLogger()
)

func CheckRelatedFile(file string) (string, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true
	}
	return "", false
}

func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)
	if !ok {
		log.Fatal("Composer Handler received unknown type")
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
		stream.AddComponent(component)
	}

	return data
}
