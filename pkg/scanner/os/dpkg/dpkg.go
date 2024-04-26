package dpkg

import (
	"slices"
	"strings"

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
	payload, ok := data.(types.Payload)
	if !ok {
		log.Error("Dpkg received unknown file type")
		return nil
	}

	scan(payload)

	return data
}

func scan(payload types.Payload) {
	manifest := payload.Body.(types.ManifestFile)
	contents := string(manifest.Content)
	packages := helper.SplitContentsByEmptyLine(contents)

	for _, info := range packages {
		metadata := parseMetadata(info)

		if metadata["package"] == nil || metadata["version"] == nil {
			continue
		}

		n, ok := metadata["package"].(string)
		if !ok {
			continue
		}

		version, ok := metadata["version"].(string)
		if !ok {
			continue
		}

		name := cleanName(n)

		var desc string
		if val, ok := metadata["description"].(string); ok {
			desc = val
		}

		c := component.New(name, version, Type)

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

		rawMetadata, err := helper.ToJSON(metadata)
		if err != nil {
			log.Errorf("Error converting metadata to JSON: %s", err)
		}

		if len(rawMetadata) > 0 {
			component.AddRawMetadata(c, rawMetadata)
		}

		cdx.AddComponent(c, payload.Address)

		if metadata["source"] != nil {

			n, ok := metadata["source"].(string)
			if !ok {
				continue
			}

			version, ok := metadata["version"].(string)
			if !ok {
				continue
			}

			name := cleanName(n)

			o := component.New(name, version, Type)

			cpes := cpe.NewCPE23(o.Name, o.Name, o.Version, Type)
			if len(cpes) > 0 {
				for _, cpe := range cpes {
					component.AddCPE(o, cpe)
				}
			}

			component.AddOrigin(o, manifest.Path)
			component.AddType(o, Type)

			rawMetadata, err := helper.ToJSON(metadata)
			if err != nil {
				log.Errorf("Error converting metadata to JSON: %s", err)
			}

			if len(rawMetadata) > 0 {
				component.AddRawMetadata(o, rawMetadata)
			}

			cdx.AddComponent(o, payload.Address)
		}
	}
}

func cleanName(name string) string {
	if strings.Contains(name, "(") {
		return strings.TrimSpace(strings.Split(name, "(")[0])
	}
	return name
}
