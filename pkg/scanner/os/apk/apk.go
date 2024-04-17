package apk

import (
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
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
	payload, ok := data.(types.Payload)
	if !ok {
		log.Error("Apk Handler received unknown type")
		return nil
	}

	scan(payload)

	return data
}

func scan(payload types.Payload) {
	manifest := payload.Body.(types.ManifestFile)
	packages := strings.Split(string(manifest.Content), "\n\n")

	for _, info := range packages {
		info = strings.TrimSpace(info)
		attributes := strings.Split(info, "\n")

		metadata := parseMetadata(attributes)
		if metadata["Name"] == nil || metadata["Version"] == nil {
			continue
		}

		c := component.New(metadata["Name"].(string), metadata["Version"].(string), Type)

		cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
		if len(cpes) > 0 {
			for _, cpe := range cpes {
				component.AddCPE(c, cpe)
			}
		}

		component.AddOrigin(c, manifest.Path)
		component.AddType(c, Type)

		for _, license := range strings.Split(metadata["License"].(string), " ") {
			if strings.Contains(strings.ToLower(license), "and") {
				licenses := strings.Split(license, "and")
				for _, l := range licenses {
					component.AddLicense(c, l)
				}
			} else {
				component.AddLicense(c, license)
			}
		}

		rawMetadata, err := helper.ToJSON(metadata)
		if err != nil {
			log.Errorf("Error converting metadata to JSON: %s", err)
		}

		if len(rawMetadata) > 0 {
			component.AddRawMetadata(c, rawMetadata)
		}

		cdx.AddComponent(c, payload.Address)

		// Add origin component
		if metadata["Origin"] != nil {

			name, version := metadata["Origin"].(string), metadata["Version"].(string)

			if len(name) == 0 || len(version) == 0 {
				continue
			}

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
