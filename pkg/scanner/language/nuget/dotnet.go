package nuget

import (
	"strings"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/cdx/dependency"
	"github.com/carbonetes/diggity/pkg/types"
)

// Scans .vsproj and .vbproj files for packages
func scanProjectFile(payload types.Payload) {
	file, ok := payload.Body.(types.ManifestFile)
	if !ok {
		log.Debug("Nuget Handler received unknown type")
		return
	}

	// Read the file
	metadata, err := parseProjectFile(file.Content)
	if err != nil {
		log.Debug("Failed to parse project file")
		return
	}

	// Check if property group has package id and version information
	var dependencyNode *cyclonedx.Dependency
	c := processPropertyGroup(metadata.PropertyGroup)
	if c != nil {
		component.AddLayer(c, payload.Layer)
		component.AddOrigin(c, file.Path)
		cdx.AddComponent(c, payload.Address)
		dependencyNode = &cyclonedx.Dependency{
			Ref:          c.BOMRef,
			Dependencies: &[]string{},
		}
	}

	// Check if the file has any package references
	if len(metadata.ItemGroup.PackageReferences) != 0 {
		for _, packageReference := range metadata.ItemGroup.PackageReferences {
			if packageReference.Include != "" && (packageReference.Version != "") {
				if strings.Contains(packageReference.Version, "${") {
					continue
				}

				c := component.New(packageReference.Include, packageReference.Version, Type)
				addCPEs(c, packageReference.Include, packageReference.Version)
				rawMetadata, err := helper.ToJSON(packageReference)
				if err != nil {
					log.Debug("Failed to marshal metadata")
				}

				if len(rawMetadata) > 0 {
					component.AddRawMetadata(c, rawMetadata)
				}

				component.AddLayer(c, payload.Layer)
				component.AddOrigin(c, file.Path)
				cdx.AddComponent(c, payload.Address)

				// Add dependency node to the component
				if dependencyNode != nil {
					*dependencyNode.Dependencies = append(*dependencyNode.Dependencies, c.BOMRef)
				}
			}
		}
	}

	// Check if the file has any references
	if len(metadata.ItemGroup.References) != 0 {
		for _, reference := range metadata.ItemGroup.References {
			c := processReference(reference)
			if c != nil {
				component.AddLayer(c, payload.Layer)
				component.AddOrigin(c, file.Path)
				cdx.AddComponent(c, payload.Address)

				// Add dependency node to the component
				if dependencyNode != nil {
					*dependencyNode.Dependencies = append(*dependencyNode.Dependencies, c.BOMRef)
				}
			}
		}
	}

	if dependencyNode == nil {
		return
	}

	if len(*dependencyNode.Dependencies) > 0 {
		dependency.AddDependency(payload.Address, dependencyNode)
	}

}

func processPropertyGroup(propertyGroup PropertyGroup) *cyclonedx.Component {
	var name, version string
	if propertyGroup.PackageId != "" && (propertyGroup.PackageVersion != "" || propertyGroup.Version != "") {
		name = propertyGroup.PackageId
		if propertyGroup.PackageVersion != "" {
			version = propertyGroup.PackageVersion
		} else {
			version = propertyGroup.Version
		}

		if strings.Contains(version, "${") {
			return nil
		}
	}

	if name == "" || version == "" || strings.Contains(version, "${") {
		return nil
	}

	rawMetadata, err := helper.ToJSON(propertyGroup)
	if err != nil {
		log.Debug("Failed to marshal metadata")
	}

	c := component.New(name, version, Type)

	if len(rawMetadata) > 0 {
		component.AddRawMetadata(c, rawMetadata)
	}

	if propertyGroup.PackageRequireLicenseAcceptance {
		component.AddLicense(c, "LicenseRef-proprietary")
	}

	if propertyGroup.PackageLicenseExpression != "" {
		component.AddLicense(c, propertyGroup.PackageLicenseExpression)
	}

	addCPEs(c, name, version)

	return c
}

func processReference(ref Reference) *cyclonedx.Component {
	if ref.Include != "" && ref.Version != "" {
		c := component.New(ref.Include, ref.Version, Type)
		addCPEs(c, ref.Include, ref.Version)
		rawMetadata, err := helper.ToJSON(ref)
		if err != nil {
			log.Debug("Failed to marshal metadata")
		}

		if len(rawMetadata) > 0 {
			component.AddRawMetadata(c, rawMetadata)
		}

		return c
	}

	return nil
}