package pub

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "pub"

var (
	Manifests = []string{"pubspec.yaml", "pubspec.lock"}
	log       = logger.GetLogger()
)

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)
	if !ok {
		log.Fatal("Pub Handler received unknown type")
	}

	if strings.Contains(manifest.Path, "pubspec.yaml") {
		metadata := readManifestFile(manifest.Content)
		var name, version, license string

		name, ok = metadata["name"].(string)
		if !ok {
			return nil
		}
		version, ok = metadata["version"].(string)
		if !ok {
			version = "0.0.0"
		}

		if val, ok := metadata["license"].(string); ok {
			license = val
		}
		component := types.NewComponent(name, version, Type, manifest.Path, "", metadata)
		component.Licenses = append(component.Licenses, license)
		stream.AddComponent(component)
	} else if strings.Contains(manifest.Path, "pubspec.lock") {
		metadata := readLockFile(manifest.Content)
		if len(metadata.Packages) == 0 {
			return nil
		}
		for _, pkg := range metadata.Packages {
			if pkg.Description.Name == "" || pkg.Version == "" {
				continue
			}
			component := types.NewComponent(pkg.Description.Name, pkg.Version, Type, manifest.Path, "", pkg)
			stream.AddComponent(component)
		}
	}

	return data
}
