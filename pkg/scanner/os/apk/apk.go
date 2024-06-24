package apk

import (
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
		log.Debug("Apk Handler received unknown type")
		return nil
	}

	scan(payload)

	return data
}

func scan(payload types.Payload) {
	apkDb := payload.Body.(types.ManifestFile)
	records, err := ParseApkIndexFile(string(apkDb.Content))
	if err != nil {
		log.Debugf("error parsing apk index file: %s", err)
		return
	}
	if len(records) == 0 {
		return
	}

	for _, record := range records {
		c := component.New(record.Package, record.Version, Type)

		cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
		if len(cpes) > 0 {
			for _, cpe := range cpes {
				component.AddCPE(c, cpe)
			}
		}

		component.AddOrigin(c, apkDb.Path)
		component.AddType(c, Type)

		if record.Description != "" {
			c.Description = record.Description
		}

		if record.Maintainer != "" {
			c.Publisher = record.Maintainer
		}

		if len(record.Licenses) != 0 {
			for _, license := range record.Licenses {
				component.AddLicense(c, license)
			}
		}

		qs := map[string]string{}
		if len(record.Architecture) > 0 {
			qs["arch"] = record.Architecture
		}

		if len(record.Origin) > 0 {
			qs["upstream"] = record.Origin
		}

		component.AddRefQualifier(c, qs)

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
		if len(record.URL) > 0 {
			c.ExternalReferences = &[]cyclonedx.ExternalReference{}
			*c.ExternalReferences = append(*c.ExternalReferences, cyclonedx.ExternalReference{
				Type: cyclonedx.ERTypeDistribution,
				URL:  record.URL,
			})
		}

		dependencyNode := &cyclonedx.Dependency{
			Ref:          c.BOMRef,
			Dependencies: &[]string{},
		}

		if len(record.Dependencies) > 0 {
			for _, dep := range record.Dependencies {
				if strings.HasPrefix(dep, "so:") {
					dep = findProvider(dep, records)
				}
				dep = helper.SplitAny(dep, "=<>")[0]
				if !helper.StringSliceContains(*dependencyNode.Dependencies, dep) && !strings.Contains(dep, "/") {
					*dependencyNode.Dependencies = append(*dependencyNode.Dependencies, dep)
				}
			}

		}

		if len(*dependencyNode.Dependencies) > 0 {
			dependency.AddDependency(payload.Address, dependencyNode)
		}

		cdx.AddComponent(c, payload.Address)
	}

}

func findProvider(so string, records []ApkIndexRecord) string {
	for _, record := range records {
		for _, provides := range record.Provides {
			if strings.Contains(provides, so) {
				return record.Package
			}
		}
	}
	return so
}
