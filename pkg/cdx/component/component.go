package component

import (
	"fmt"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/google/uuid"
)

// New creates a new cyclonedx.Component with the given name, version, and category.
func New(name, version, category string) *cyclonedx.Component {
	purl := fmt.Sprintf("pkg:%s/%s@%s", category, name, version)
	return &cyclonedx.Component{
		Type:       cyclonedx.ComponentTypeLibrary,
		BOMRef:     uuid.New().String(),
		Name:       helper.CleanValue(name).(string),
		Version:    helper.CleanValue(version).(string),
		PackageURL: helper.CleanValue(purl).(string),
		Properties: &[]cyclonedx.Property{},
	}
}

// AddCPE adds a CPE to the given cyclonedx.Component.
// The CPE should be a CPE 2.3 identifier.
func AddCPE(c *cyclonedx.Component, cpe string) {
	if c.Properties == nil {
		c.Properties = &[]cyclonedx.Property{}
	}

	if len(cpe) == 0 {
		return
	}

	v := helper.CleanValue(cpe)

	*c.Properties = append(*c.Properties, cyclonedx.Property{
		Name:  "diggity:package:cpe23",
		Value: v.(string),
	})
}

// AddOrigin adds an origin to the given cyclonedx.Component.
// The origin should be the package's location on the filesystem.
func AddOrigin(c *cyclonedx.Component, origin string) {
	if c.Properties == nil {
		c.Properties = &[]cyclonedx.Property{}
	}

	if len(origin) == 0 {
		return
	}

	*c.Properties = append(*c.Properties, cyclonedx.Property{
		Name:  "diggity:file:location",
		Value: origin,
	})

}

// AddDescription adds a description to the given cyclonedx.Component.
// The description should be a found on the package's website or in the package's metadata.
func AddDescription(c *cyclonedx.Component, description string) {
	if c.Properties == nil {
		c.Properties = &[]cyclonedx.Property{}
	}

	if len(description) == 0 {
		return
	}

	v := helper.CleanValue(description)

	*c.Properties = append(*c.Properties, cyclonedx.Property{
		Name:  "diggity:package:description",
		Value: v.(string),
	})
}

// AddType adds a type to the given cyclonedx.Component.
// The type should be one of the package types defined in the Scanner Module.
func AddType(c *cyclonedx.Component, componentType string) {
	if c.Properties == nil {
		c.Properties = &[]cyclonedx.Property{}
	}

	if len(componentType) == 0 {
		return
	}

	*c.Properties = append(*c.Properties, cyclonedx.Property{
		Name:  "diggity:package:type",
		Value: componentType,
	})
}

// AddRawMetadata adds raw metadata to the given cyclonedx.Component.
// The metadata should be in string value.
func AddRawMetadata(c *cyclonedx.Component, metadata []byte) {
	if c.Properties == nil {
		c.Properties = &[]cyclonedx.Property{}
	}

	if len(metadata) == 0 {
		return
	}

	v, err := helper.CleanJSON(string(metadata))
	if err != nil {
		log.Error(err)
	}

	*c.Properties = append(*c.Properties, cyclonedx.Property{
		Name:  "diggity:package:metadata",
		Value: v,
	})
}

// AddLicense adds a license to the given cyclonedx.Component.
// The license should be a SPDX license identifier.
//
//	https://spdx.org/licenses/
func AddLicense(c *cyclonedx.Component, license string) {
	if c.Licenses == nil {
		c.Licenses = &cyclonedx.Licenses{}
	}

	if len(license) == 0 {
		return
	}

	v := helper.CleanValue(license)

	*c.Licenses = append(*c.Licenses, cyclonedx.LicenseChoice{
		License: &cyclonedx.License{
			ID: v.(string),
		},
	})
}
