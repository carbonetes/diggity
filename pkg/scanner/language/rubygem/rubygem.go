package rubygem

import (
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "gem"

var Manifests = []string{"Gemfile.lock"}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) || (strings.Contains(file, ".gemspec") && strings.Contains(file, "specifications")) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)
	if !ok {
		log.Error("Rubygem Handler received unknown type")
	}

	if strings.Contains(manifest.Path, "Gemfile.lock") {
		attributes := readManifestFile(manifest.Content)
		for _, attribute := range attributes {
			name, version := attribute[0], attribute[1]
			metadata := map[string]string{"name": name, "version": version}
			component := types.NewComponent(name, version, Type, manifest.Path, "", metadata)
			cpes := cpe.NewCPE23(component.Name, component.Name, component.Version, Type)
			if len(cpes) > 0 {
				component.CPEs = append(component.CPEs, cpes...)
			}
			stream.AddComponent(component)
		}
	} else if strings.Contains(manifest.Path, ".gemspec") && strings.Contains(manifest.Path, "specifications") {
		metadata := readGemspecFile(manifest.Content)
		if _, ok := metadata["metadata"].(string); ok {
			delete(metadata, "metadata")
		}
		name, version := metadata["name"].(string), metadata["version"].(string)
		component := types.NewComponent(name, version, Type, manifest.Path, "", metadata)
		if val, ok := metadata["description"].(string); ok {
			component.Description = val
		}
		var licenses []string
		if val, ok := metadata["licenses"].(string); ok {
			license := regexp.MustCompile(`[^\w^,^ ]`).ReplaceAllString(val, "")
			component.Licenses = append(licenses, license)
		}
		cpes := cpe.NewCPE23(component.Name, component.Name, component.Version, Type)
		if len(cpes) > 0 {
			component.CPEs = append(component.CPEs, cpes...)
		}
		stream.AddComponent(component)
	}

	return data
}
