package cran

import (
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "cran"

var (
	RelatedPath = "usr/lib/R/"
	RelatedFile = "DESCRIPTION"
)

func CheckRelatedFiles(file string) (string, bool, bool) {
	if strings.Contains(file, RelatedPath) && RelatedFile == filepath.Base(file) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)

	if !ok {
		log.Error("Cran Handler received unknown type")
		return nil
	}

	metadata := readManifestFile(manifest.Content)

	if metadata.Package == "" || metadata.Version == "" {
		return nil
	}

	component := types.NewComponent(metadata.Package, metadata.Version, Type, manifest.Path, metadata.Description, metadata)
	stream.AddComponent(component)
	return data
}
