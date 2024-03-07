package component

import (
	"fmt"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/google/uuid"
)

// New creates a new cyclonedx.Component with the given name, version, and category.
func New(name, version, category string) *cyclonedx.Component {
	return &cyclonedx.Component{
		Type:       cyclonedx.ComponentTypeLibrary,
		BOMRef:     uuid.New().String(),
		Name:       name,
		Version:    version,
		PackageURL: fmt.Sprintf("pkg:%s/%s@%s", category, name, version),
		Properties: &[]cyclonedx.Property{},
	}
}

// AddCPE adds a CPE to the given cyclonedx.Component.
// The CPE should be a CPE 2.3 identifier.
func AddCPE(c *cyclonedx.Component, cpe string) {
	if c.Properties == nil {
		c.Properties = &[]cyclonedx.Property{}
	}

	*c.Properties = append(*c.Properties, cyclonedx.Property{
		Name:  "cpe23",
		Value: cpe,
	})
}

// AddOrigin adds an origin to the given cyclonedx.Component.
// The origin should be the package's location on the filesystem.
func AddOrigin(c *cyclonedx.Component, origin string) {
	if c.Properties == nil {
		c.Properties = &[]cyclonedx.Property{}
	}

	*c.Properties = append(*c.Properties, cyclonedx.Property{
		Name:  "location",
		Value: origin,
	})

}

// AddDescription adds a description to the given cyclonedx.Component.
// The description should be a found on the package's website or in the package's metadata.
func AddDescription(c *cyclonedx.Component, description string) {
	if c.Properties == nil {
		c.Properties = &[]cyclonedx.Property{}
	}

	*c.Properties = append(*c.Properties, cyclonedx.Property{
		Name:  "diggity:package:description",
		Value: description,
	})
}

// AddType adds a type to the given cyclonedx.Component.
// The type should be one of the package types defined in the Scanner Module.
func AddType(c *cyclonedx.Component, componentType string) {
	if c.Properties == nil {
		c.Properties = &[]cyclonedx.Property{}
	}

	*c.Properties = append(*c.Properties, cyclonedx.Property{
		Name:  "diggity:package:type",
		Value: componentType,
	})
}

// AddLicense adds a license to the given cyclonedx.Component.
// The license should be a SPDX license identifier.
//	https://spdx.org/licenses/
func AddLicense(c *cyclonedx.Component, license string) {
	if c.Licenses == nil {
		c.Licenses = &cyclonedx.Licenses{}
	}

	*c.Licenses = append(*c.Licenses, cyclonedx.LicenseChoice{
		License: &cyclonedx.License{
			ID: license,
		},
	})
}
