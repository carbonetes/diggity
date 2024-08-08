package dpkg

import (
	"path/filepath"
	"strings"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/cdx/dependency"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "deb"

var (
	RelatedPath = "dpkg/"
	RelatedFile = "status"
)

func CheckRelatedFile(file string) (string, bool, bool) {
	if strings.Contains(file, RelatedPath) {
		if filepath.Base(file) == RelatedFile {
			return Type, true, true
		}
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	payload, ok := data.(types.Payload)
	if !ok {
		log.Debug("Dpkg received unknown file type")
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

	contents := string(file.Content)

	records, err := ParseDpkgDatabase(contents)
	if err != nil {
		log.Debugf("error parsing dpkg database: %s", err)
		return
	}

	for _, record := range records {
		if record.Name == "" || record.Version == "" {
			continue
		}

		name := cleanName(record.Name)

		c := component.New(name, record.Version, Type)

		cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
		if len(cpes) > 0 {
			for _, cpe := range cpes {
				component.AddCPE(c, cpe)
			}
		}

		component.AddOrigin(c, file.Path)
		component.AddType(c, Type)

		if record.Description != "" {
			component.AddDescription(c, record.Description)
		}

		rawMetadata, err := helper.ToJSON(record)
		if err != nil {
			log.Debugf("Error converting metadata to JSON: %s", err)
		}

		if len(rawMetadata) > 0 {
			component.AddRawMetadata(c, rawMetadata)
		}

		if len(payload.Layer) > 0 {
			component.AddLayer(c, payload.Layer)
		}

		if record.Maintainer != "" {
			c.Publisher = record.Maintainer
		}

		if record.Homepage != "" {
			c.ExternalReferences = &[]cyclonedx.ExternalReference{
				{
					Type: "website",
					URL:  record.Homepage,
				},
			}
		}

		qm := make(map[string]string)
		if record.Architecture != "" {
			qm["arch"] = record.Architecture
		}

		if record.Source != "" {
			qm["upstream"] = record.Source
		}

		component.AddRefQualifier(c, qm)

		dependencyNode := &cyclonedx.Dependency{
			Ref:          c.BOMRef,
			Dependencies: &[]string{},
		}

		if len(record.Depends) > 0 {
			for _, entry := range record.Depends {
				// *dependencyNode.Dependencies = append(*dependencyNode.Dependencies, entry...)
				for _, dep := range entry {
					*dependencyNode.Dependencies = append(*dependencyNode.Dependencies, findProvider(dep, records))
				}
			}
		}

		if len(record.PreDepends) > 0 {
			for _, entry := range record.PreDepends {
				for _, dep := range entry {
					dependency := findProvider(dep, records)
					if dependency != "" {
						*dependencyNode.Dependencies = append(*dependencyNode.Dependencies, dependency)
					}
				}
			}
		}

		if len(*dependencyNode.Dependencies) > 0 {
			dependency.AddDependency(payload.Address, dependencyNode)
		}

		cdx.AddComponent(c, payload.Address)

	}

}

func cleanName(name string) string {
	if strings.Contains(name, "(") {
		return strings.TrimSpace(strings.Split(name, "(")[0])
	}
	return name
}

func findProvider(dep string, records []Package) string {
	for _, record := range records {
		if record.Provides == nil {
			continue
		}
		for _, provider := range record.Provides {
			if strings.Contains(provider, dep) {
				return record.Name
			}
		}
	}
	return dep
}
