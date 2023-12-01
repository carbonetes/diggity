package nuget

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "nuget"

var (
	Manifests = []string{".deps.json"}
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
		log.Error("Nuget Handler received unknown type")
	}

	metadata := readManifestFile(manifest.Content)
	if len(metadata.Libraries) == 0 {
		return nil
	}

	for id, pkg := range metadata.Libraries {
		if pkg.Type != "package" {
			continue
		}

		attributes := strings.Split(id, "/")
		if len(attributes) != 2 {
			continue
		}

		name, version := attributes[0], attributes[1]

		component := types.NewComponent(name, version, Type, manifest.Path, "", pkg)
		stream.AddComponent(component)
	}
	return data
}
