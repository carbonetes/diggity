package pacman

import (
	"path/filepath"
	"slices"

	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "archlinux"

var (
	Manifests = []string{"var/lib/pacman/local"}
	log       = logger.GetLogger()
)

func CheckRelatedFile(file string) (string, bool) {
	if slices.Contains(Manifests, file) && filepath.Base(file) == "desc" {
		return Type, true
	}
	return "", false
}

func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)
	if !ok {
		log.Error("Pacman Handler received unknown type")
		return nil
	}

	contents := string(manifest.Content)
	attributes := helper.SplitContentsByEmptyLine(contents)
	metadata := parseMetadata(attributes)

	if metadata["name"] == nil || metadata["name"] == "" {
		return nil
	}

	name, version, desc := metadata["name"].(string), metadata["version"].(string), metadata["description"].(string)
	component := types.NewComponent(name, version, Type, manifest.Path, desc, metadata)

	arch, ok := metadata["arch"].(string)
	if !ok {
		arch = ""
	}
	component.PURL = component.PURL + "?arch=" + arch
	stream.AddComponent(component)

	return data
}
