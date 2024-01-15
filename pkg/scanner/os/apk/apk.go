package apk

import (
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "apk"

var RelatedPath = "lib/apk/db/installed"

func CheckRelatedFile(file string) (string, bool, bool) {
	if strings.Contains(file, RelatedPath) {

		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)
	if !ok {
		log.Error("Apk Handler received unknown type")
		return nil
	}

	packages := strings.Split(string(manifest.Content), "\n\n")

	for _, info := range packages {
		info = strings.TrimSpace(info)
		attributes := strings.Split(info, "\n")

		metadata := parseMetadata(attributes)
		if metadata["Name"] == nil || metadata["Version"] == nil {
			continue
		}

		component := types.NewComponent(metadata["Name"].(string), metadata["Version"].(string), Type, manifest.Path, metadata["Description"].(string), metadata)
		for _, license := range strings.Split(metadata["License"].(string), " ") {
			if !strings.Contains(strings.ToLower(license), "and") {
				component.Licenses = append(component.Licenses, license)
			}
		}
		cpes := cpe.NewCPE23(component.Name, component.Name, component.Version, Type)
		if len(cpes) > 0 {
			component.CPEs = append(component.CPEs, cpes...)
		}
		stream.AddComponent(component)

		if metadata["Origin"] != nil {
			origin := types.NewComponent(metadata["Origin"].(string), metadata["Version"].(string), Type, manifest.Path, metadata["Description"].(string), nil)
			origin.Licenses = component.Licenses
			cpes := cpe.NewCPE23(component.Name, component.Name, component.Version, Type)
			if len(cpes) > 0 {
				origin.CPEs = append(component.CPEs, cpes...)
			}
			stream.AddComponent(origin)
		}
	}
	return data
}
