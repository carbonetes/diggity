package rubygem

import (
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
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

			c := component.New(name, version, Type)

			cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
			if len(cpes) > 0 {
				for _, cpe := range cpes {
					component.AddCPE(c, cpe)
				}
			}

			component.AddOrigin(c, manifest.Path)
			component.AddType(c, Type)

			cdx.AddComponent(c)
		}
	} else if strings.Contains(manifest.Path, ".gemspec") && strings.Contains(manifest.Path, "specifications") {
		metadata := readGemspecFile(manifest.Content)
		if _, ok := metadata["metadata"].(string); ok {
			delete(metadata, "metadata")
		}
		name, version := metadata["name"].(string), metadata["version"].(string)

		c := component.New(name, version, Type)

		cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
		if len(cpes) > 0 {
			for _, cpe := range cpes {
				component.AddCPE(c, cpe)
			}
		}

		component.AddOrigin(c, manifest.Path)
		component.AddType(c, Type)

		if val, ok := metadata["licenses"].(string); ok {
			license := regexp.MustCompile(`[^\w^,^ ]`).ReplaceAllString(val, "")
			component.AddLicense(c, license)
		}

		if val, ok := metadata["description"].(string); ok {
			component.AddDescription(c, val)
		}

		cdx.AddComponent(c)
	}

	return data
}
