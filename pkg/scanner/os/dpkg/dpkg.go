package dpkg

import (
	"slices"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
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

		c := component.New(metadata["package"].(string), metadata["version"].(string), Type)

		cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
		if len(cpes) > 0 {
			for _, cpe := range cpes {
				component.AddCPE(c, cpe)
			}
		}

		component.AddOrigin(c, manifest.Path)
		component.AddType(c, Type)

		if desc != "" {
			component.AddDescription(c, desc)
		}

		cdx.AddComponent(c)

		if metadata["source"] != nil {

			o := component.New(metadata["source"].(string), metadata["version"].(string), Type)

			cpes := cpe.NewCPE23(o.Name, o.Name, o.Version, Type)
			if len(cpes) > 0 {
				for _, cpe := range cpes {
					component.AddCPE(o, cpe)
				}
			}

			component.AddOrigin(o, manifest.Path)
			component.AddType(o, Type)

			cdx.AddComponent(o)
		}
	}

	return data
}
