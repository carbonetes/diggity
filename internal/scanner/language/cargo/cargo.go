package cargo

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "rust-crate"

var (
	Manifests = []string{"Cargo.lock"}
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
		log.Fatal("Cargo Handler received unknown type")
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
		stream.AddComponent(component)
	}

	return data
}
