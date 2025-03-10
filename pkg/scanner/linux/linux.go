package linux

import (
	"fmt"
	"slices"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/types"
)

var (
	Releases []types.OSRelease
	Type     = "osrelease"
	// Add more os release files here if needed
	Manifests      = []string{"etc/os-release", "usr/lib/os-release", "etc/lsb-release", "etc/centos-release", "etc/redhat-release", "etc/debian_version", "etc/alpine-release", "etc/SuSE-release", "etc/gentoo-release", "etc/arch-release", "etc/oracle-release"}
	PropertyPrefix = "diggity:" + Type + ":"
)

func Scan(data interface{}) interface{} {
	payload, ok := data.(types.Payload)
	if !ok {
		log.Debug("Distro handler received unknown type")
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

	Releases = append(Releases, parse(file))

	for _, release := range Releases {
		var name, version, desc string
		if val, ok := release.Release["id"].(string); ok {
			name = val
		}

		if val, ok := release.Release["version_id"].(string); ok {
			version = val
		}

		if val, ok := release.Release["pretty_name"].(string); ok {
			desc = val
		}

		c := newOSComponent(name, version, desc)

		if c == nil {
			continue
		}

		swid := cyclonedx.SWID{
			TagID:   name,
			Name:    name,
			Version: version,
		}

		for key, value := range release.Release {
			if key == "home_url" {
				addExternalReference(c, value.(string), "website")
				swid.URL = value.(string)
				continue
			}

			if key == "support_url" {
				addExternalReference(c, value.(string), "support")
				continue
			}

			if key == "bug_report_url" {
				addExternalReference(c, value.(string), "issue-tracker")
				continue
			}

			if key == "privacy_policy_url" {
				addExternalReference(c, value.(string), "privacy-policy")
				continue
			}

			if key == "cpe_name" {
				addProperty(c, PropertyPrefix+"cpe", value.(string))
				continue
			}

			if key == "documentation_url" {
				addExternalReference(c, value.(string), "documentation")
				continue
			}

			addProperty(c, PropertyPrefix+key, value.(string))
		}

		c.SWID = &swid
		addProperty(c, PropertyPrefix+"location", release.File)

		cdx.AddComponent(c, payload.Address)
	}
}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, file) {
		return Type, true, true
	}
	return "", false, false
}

func newOSComponent(name, version, desc string) *cyclonedx.Component {
	if name == "" && version == "" {
		return nil
	}

	c := &cyclonedx.Component{
		Type:        cyclonedx.ComponentTypeOS,
		BOMRef:      fmt.Sprintf("os:%s@%s", name, version),
		Name:        name,
		Version:     version,
		Description: desc,
	}

	return c
}

func addProperty(c *cyclonedx.Component, name, value string) {
	if c.Properties == nil {
		c.Properties = &[]cyclonedx.Property{}
	}

	*c.Properties = append(*c.Properties, cyclonedx.Property{
		Name:  name,
		Value: value,
	})
}

func addExternalReference(c *cyclonedx.Component, url, desc string) {
	if c.ExternalReferences == nil {
		c.ExternalReferences = &[]cyclonedx.ExternalReference{}
	}

	*c.ExternalReferences = append(*c.ExternalReferences, cyclonedx.ExternalReference{
		URL:  url,
		Type: cyclonedx.ExternalReferenceType(desc),
	})
}
