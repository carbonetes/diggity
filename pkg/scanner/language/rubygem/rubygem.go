package rubygem

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/helper"
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
	payload, ok := data.(types.Payload)
	if !ok {
		log.Debug("Rubygem Handler received unknown type")
		return nil
	}

	scan(payload)

	return data
}

func scan(payload types.Payload) {
	file, ok := payload.Body.(types.ManifestFile)
	if !ok {
		log.Debugf("Failed to convert payload body to manifest file")
		return
	}

	if strings.Contains(file.Path, "Gemfile.lock") {
		attributes := readManifestFile(file.Content)
		for _, attribute := range attributes {
			name, version := attribute[0], attribute[1]

			c := component.New(name, version, Type)

			cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
			if len(cpes) > 0 {
				for _, cpe := range cpes {
					component.AddCPE(c, cpe)
				}
			}

			component.AddOrigin(c, file.Path)
			component.AddType(c, Type)

			rawMetadata, err := helper.ToJSON(attribute)
			if err != nil {
				log.Debugf("Error converting metadata to JSON: %s", err)
			}

			if len(rawMetadata) > 0 {
				component.AddRawMetadata(c, rawMetadata)
			}

			if len(payload.Layer) > 0 {
				component.AddLayer(c, payload.Layer)
			}

			cdx.AddComponent(c, payload.Address)
		}
	} else if strings.Contains(file.Path, ".gemspec") && strings.Contains(file.Path, "specifications") {
		metadata, err := parseGemspec(file.Content)
		if err != nil {
			log.Debugf("error parsing gemspec: %s", err)
			return
		}

		if metadata == nil {
			return
		}

		cleanMetadata(metadata)

		if metadata.Name == "" || metadata.Version == "" {
			return
		}

		name, version := metadata.Name, metadata.Version

		c := component.New(name, version, Type)

		cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
		if len(cpes) > 0 {
			for _, cpe := range cpes {
				component.AddCPE(c, cpe)
			}
		}

		component.AddOrigin(c, file.Path)
		component.AddType(c, Type)

		if len(metadata.Authors) > 0 {
			c.Authors = &[]cyclonedx.OrganizationalContact{}
			for _, author := range metadata.Authors {
				*c.Authors = append(*c.Authors, cyclonedx.OrganizationalContact{Name: author})
			}
		}

		if metadata.Description != "" {
			c.Description = metadata.Description
		}

		if len(metadata.Licenses) > 0 {
			for _, license := range metadata.Licenses {
				component.AddLicense(c, license)
			}
		}

		rawMetadata, err := helper.ToJSON(metadata)
		if err != nil {
			log.Debugf("Error converting metadata to JSON: %s", err)
		}

		if len(rawMetadata) > 0 {
			component.AddRawMetadata(c, rawMetadata)
		}

		if len(payload.Layer) > 0 {
			component.AddLayer(c, payload.Layer)
		}

		cdx.AddComponent(c, payload.Address)
	}
}
