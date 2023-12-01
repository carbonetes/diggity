package hex

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "hex"

var (
	Manifests = []string{"rebar.lock", "mix.lock"}
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
		log.Error("Hex Handler received unknown type")
	}

	if strings.Contains(manifest.Path, "rebar.lock") {
		packages := readRebarFile(manifest.Content)
		if len(packages) == 0 {
			return nil
		}

		for _, pkg := range packages {
			if pkg.Name == "" || pkg.Version == "" {
				continue
			}

			component := types.NewComponent(pkg.Name, pkg.Version, Type, manifest.Path, "", pkg)
			stream.AddComponent(component)
		}

	} else if strings.Contains(manifest.Path, "mix.lock") {
		packages := readMixFile(manifest.Content)
		if len(packages) == 0 {
			return nil
		}

		for _, pkg := range packages {
			if pkg.Name == "" || pkg.Version == "" {
				continue
			}

			component := types.NewComponent(pkg.Name, pkg.Version, Type, manifest.Path, "", pkg)
			stream.AddComponent(component)
		}
	}

	return data
}
