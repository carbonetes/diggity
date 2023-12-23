package convert

import (
	"encoding/xml"
	"fmt"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/pkg/types"
)

var (
	// XMLN cyclonedx
	XMLN = fmt.Sprintf("http://cyclonedx.org/schema/bom/%+v", cyclonedx.SpecVersion1_5)
)

const (
	cycloneDX        = "CycloneDX"
	vendor           = "carbonetes"
	name             = "diggity"
	diggityPrefix    = "diggity"
	packagePrefix    = "package"
	distroPrefix     = "distro"
	colonPrefix      = ":"
	cpePrefix        = "cpe23"
	locationPrefix   = "location"
	library          = "library"
	operatingSystem  = "operating-system"
	issueTracker     = "issue-tracker"
	referenceWebsite = "website"
	referenceOther   = "other"
	version          = 1
)

func ToCDX(sbom *types.SBOM) *cyclonedx.BOM {
	bom := &cyclonedx.BOM{
		XMLName:      xml.Name{Local: cycloneDX},
		XMLNS:        XMLN,
		BOMFormat:    cycloneDX,
		Version:      version,
		SerialNumber: sbom.Serial,
		SpecVersion:  cyclonedx.SpecVersion1_5,
		Metadata: &cyclonedx.Metadata{
			Tools: &[]cyclonedx.Tool{},
		},
		Components: &[]cyclonedx.Component{},
	}

	bom.Metadata.Timestamp = sbom.Timestamp.String()
	*bom.Metadata.Tools = append(*bom.Metadata.Tools, cyclonedx.Tool{
		Vendor:  vendor,
		Name:    name,
		Version: sbom.Version,
	})

	for _, component := range sbom.Components {
		*bom.Components = append(*bom.Components, *ToCDXComponent(&component))
	}

	return bom
}

func ToCDXComponent(component *types.Component) *cyclonedx.Component {
	var licenses []cyclonedx.LicenseChoice

	for _, license := range component.Licenses {
		licenses = append(licenses, cyclonedx.LicenseChoice{
			License: &cyclonedx.License{
				ID: license,
			},
		})
	}

	c := &cyclonedx.Component{
		Type:       library,
		BOMRef:     component.ID,
		Name:       component.Name,
		Version:    component.Version,
		PackageURL: component.PURL,
		Properties: &[]cyclonedx.Property{},
	}

	if len(licenses) > 0 {
		*c.Licenses = append(*c.Licenses, licenses...)
	}

	if component.CPEs != nil {
		for _, cpe := range component.CPEs {
			*c.Properties = append(*c.Properties, cyclonedx.Property{
				Name:  cpePrefix,
				Value: cpe,
			})
		}
	}

	if component.Origin != "" {
		*c.Properties = append(*c.Properties, cyclonedx.Property{
			Name:  locationPrefix ,
			Value: component.Origin,
		})
	}

	if component.Description != "" {
		*c.Properties = append(*c.Properties, cyclonedx.Property{
			Name:  diggityPrefix + colonPrefix + packagePrefix + colonPrefix + "description",
			Value: component.Description,
		})
	}

	if component.Type != "" {
		*c.Properties = append(*c.Properties, cyclonedx.Property{
			Name:  diggityPrefix + colonPrefix + packagePrefix + colonPrefix + "type",
			Value: component.Type,
		})
	}

	return c
}
