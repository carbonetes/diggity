package cargo

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "rust-crate"

var (
	Manifests = []string{"Cargo.lock"}
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
		log.Error("Cargo Handler received unknown type")
	}

	packages := readManifestFile(manifest.Content)

	for _, pkg := range packages {
		if !strings.Contains(pkg, "[[package]]") {
			continue
		}

		metadata := parseMetadata(pkg)
		if metadata == nil {
			continue
		}

		if metadata["Name"] == nil {
			return nil
		}

		component := types.NewComponent(metadata["Name"].(string), metadata["Version"].(string), Type, manifest.Path, "", metadata)

		cpes := cpe.NewCPE23(component.Name, component.Name, component.Version, Type)
		if len(cpes) > 0 {
			component.CPEs = append(component.CPEs, cpes...)
		}

		stream.AddComponent(component)
	}

	return data
}
