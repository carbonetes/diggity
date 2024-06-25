package maven

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
	"github.com/carbonetes/diggity/pkg/cdx/dependency"
	"github.com/carbonetes/diggity/pkg/types"
	"github.com/golistic/urn"
)

const Type string = "java"

var Manifests = []string{"pom.xml", "pom.properties", "MANIFEST.MF"}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	payload, ok := data.(types.Payload)
	if !ok {
		log.Debug("Java Archive received unknown file type")
		return nil
	}

	scan(payload)

	return data
}

/*
TODO: Implement the following functions
 1. Implement the function to parse all manifest files in the Maven project (done)
 2. Chart out the dependencies of the Maven project
 3. Scan vendor information from the manifest files and add generate CPEs for the components
*/

func scan(payload types.Payload) {
	manifest := payload.Body.(types.ManifestFile)

	switch filepath.Base(manifest.Path) {
	case "pom.xml":
		readPOMFile(manifest, payload.Address)
	case "MANIFEST.MF":
		// readManifestFile(manifest, payload.Address) // Temporary disabled
	case "pom.properties":
		readPOMPropertiesFile(manifest, payload.Layer, payload.Address)
	}
}

//nolint:all
func readPOMFile(manifest types.ManifestFile, addr *urn.URN) {
	metadata, err := parsePOM(manifest.Content)
	if err != nil {
		log.Debugf("Failed to parse POM file: %v", err)
		return
	}

	if metadata == nil {
		return
	}

	properties := make(map[string]string)
	// Check if there are properties in the pom file
	if metadata.Properties != nil && len(metadata.Properties.Properties) > 0 {
		properties = resolveProperties(metadata)
	}

	if metadata.ArtifactID == "" || metadata.Version == "" {
		return
	}

	c := component.New(metadata.ArtifactID, metadata.Version, Type)

	cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
	if len(cpes) > 0 {
		for _, cpe := range cpes {
			component.AddCPE(c, cpe)
		}
	}

	// Correction for PackageURL
	c.PackageURL = "pkg:maven/" + metadata.GroupID + "/" + metadata.ArtifactID + "@" + metadata.Version

	component.AddOrigin(c, manifest.Path)
	component.AddType(c, Type)
	component.AddDescription(c, metadata.Description)

	meta := Dependency{
		GroupID:    metadata.GroupID,
		ArtifactID: metadata.ArtifactID,
		Version:    metadata.Version,
	}

	rawMetadata, err := helper.ToJSON(meta)
	if err != nil {
		log.Debugf("Failed to convert metadata to JSON: %v", err)
	}

	if len(rawMetadata) > 0 {
		component.AddRawMetadata(c, rawMetadata)
	}

	cdx.AddComponent(c, addr)

	if len(metadata.Dependencies) > 0 {
		dependencyNode := &cyclonedx.Dependency{
			Ref:          c.BOMRef,
			Dependencies: &[]string{},
		}

		for _, dependency := range metadata.Dependencies {

			if dependency.Version == "" {
				if (metadata.Parent.GroupID != "" && metadata.Parent.Version != "") && (dependency.GroupID == metadata.Parent.GroupID) {
					dependency.Version = metadata.Parent.Version
				}
			}

			if strings.Contains(dependency.Version, "${") {
				dependency.Version = properties[dependency.Version]
			}

			if strings.Contains(dependency.GroupID, "${") {
				dependency.GroupID = properties[dependency.GroupID]
			}

			if strings.Contains(dependency.ArtifactID, "${") {
				dependency.ArtifactID = properties[dependency.ArtifactID]
			}

			if dependency.ArtifactID == "" || dependency.Version == "" {
				continue
			}

			// Skip unresolved properties
			if strings.Contains(dependency.Version, "${") || strings.Contains(dependency.ArtifactID, "${") {
				continue
			}

			c := component.New(dependency.ArtifactID, dependency.Version, Type)

			cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
			if len(cpes) > 0 {
				for _, cpe := range cpes {
					component.AddCPE(c, cpe)
				}
			}

			// Correction for PackageURL
			c.PackageURL = "pkg:maven/" + dependency.GroupID + "/" + dependency.ArtifactID + "@" + dependency.Version
			c.BOMRef = c.PackageURL

			component.AddOrigin(c, manifest.Path)
			component.AddType(c, Type)

			if len(dependency.GroupID) > 0 {
				c.Group = dependency.GroupID
			}

			rawMetadata, err := helper.ToJSON(dependency)
			if err != nil {
				log.Debugf("Failed to convert metadata to JSON: %v", err)
			}

			if len(rawMetadata) > 0 {
				component.AddRawMetadata(c, rawMetadata)
			}

			*dependencyNode.Dependencies = append(*dependencyNode.Dependencies, dependency.ArtifactID)

			cdx.AddComponent(c, addr)
		}
		if len(*dependencyNode.Dependencies) > 0 {
			dependency.AddDependency(addr, dependencyNode)
		}
	}

}

//nolint:all
func readManifestFile(manifest types.ManifestFile, addr *urn.URN) {
	metadata, err := parseManifestFile(manifest.Content)
	if err != nil {
		log.Debugf("Failed to parse manifest file: %v", err)
		return
	}

	if len(metadata["Bundle-SymbolicName"]) == 0 || len(metadata["Bundle-Version"]) == 0 {
		return
	}

	c := component.New(metadata["Bundle-SymbolicName"], metadata["Bundle-Version"], Type)

	cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
	if len(cpes) > 0 {
		for _, cpe := range cpes {
			component.AddCPE(c, cpe)
		}
	}

	// Correction for PackageURL
	c.PackageURL = "pkg:maven/" + c.Name + "@" + c.Version

	component.AddOrigin(c, manifest.Path)
	component.AddType(c, Type)

	rawMetadata, err := helper.ToJSON(metadata)
	if err != nil {
		log.Debugf("Failed to convert metadata to JSON: %v", err)
	}

	if len(rawMetadata) > 0 {
		component.AddRawMetadata(c, rawMetadata)
	}

	cdx.AddComponent(c, addr)

}

func readPOMPropertiesFile(manifest types.ManifestFile, layer string, addr *urn.URN) {
	metadata, err := parsePOMProperties(manifest.Content)
	if err != nil {
		log.Debugf("Failed to parse POM properties file: %v", err)
		return
	}

	if metadata == nil {
		return
	}

	if len(metadata["artifactId"]) == 0 || len(metadata["version"]) == 0 {
		return
	}

	c := component.New(metadata["artifactId"], metadata["version"], Type)

	cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
	if len(cpes) > 0 {
		for _, cpe := range cpes {
			component.AddCPE(c, cpe)
		}
	}

	// Correction for PackageURL
	c.PackageURL = "pkg:maven/" + c.Name + "@" + c.Version

	component.AddOrigin(c, manifest.Path)
	component.AddType(c, Type)

	rawMetadata, err := helper.ToJSON(metadata)
	if err != nil {
		log.Debugf("Failed to convert metadata to JSON: %v", err)
	}

	if len(rawMetadata) > 0 {
		component.AddRawMetadata(c, rawMetadata)
	}

	if len(layer) > 0 {
		component.AddLayer(c, layer)
	}

	cdx.AddComponent(c, addr)

}
