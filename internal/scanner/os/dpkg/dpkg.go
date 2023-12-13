package dpkg

import (
	"slices"

	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "deb"

var Manifests = []string{"var/lib/dpkg/status"}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, file) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)

	if !ok {
		log.Error("Dpkg received unknown file type")
		return nil
	}

	contents := string(manifest.Content)
	packages := helper.SplitContentsByEmptyLine(contents)

	for _, info := range packages {
		metadata := parseMetadata(info)
		if metadata["package"] == nil || metadata["version"] == nil {
			continue
		}

		var desc string
		if val, ok := metadata["description"].(string); ok {
			desc = val
		}

		component := types.NewComponent(metadata["package"].(string), metadata["version"].(string), Type, manifest.Path, desc, metadata)
		stream.AddComponent(component)

		if metadata["source"] != nil {
			origin := types.NewComponent(metadata["source"].(string), metadata["version"].(string), Type, manifest.Path, desc, nil)
			origin.Licenses = component.Licenses
			stream.AddComponent(origin)
		}
	}

	return data
}
