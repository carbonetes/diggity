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
const ErrConvertMetadataToJSON = "Failed to convert metadata to JSON: %v"
const MavenPackageURLPrefix = "pkg:maven/"

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
	file, ok := payload.Body.(types.ManifestFile)
	if !ok {
		log.Debugf("Failed to convert payload body to manifest file")
		return
	}

	switch filepath.Base(file.Path) {
	case "pom.xml":
		readPOMFile(file, payload.Address)
	case "MANIFEST.MF":
		// readManifestFile(manifest, payload.Address) // Temporary disabled
	case "pom.properties":
		readPOMPropertiesFile(file, payload.Layer, payload.Address)
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

	properties := getProperties(metadata)

	if metadata.ArtifactID == "" || metadata.Version == "" {
		return
	}

	c := createComponent(metadata, manifest.Path)

	addCPEs(c)
	c.PackageURL = MavenPackageURLPrefix + metadata.GroupID + "/" + metadata.ArtifactID + "@" + metadata.Version

	addComponentDetails(c, manifest.Path, metadata.Description, metadata)

	cdx.AddComponent(c, addr)

	if len(metadata.Dependencies) > 0 {
		processDependencies(metadata, properties, manifest.Path, addr)
	}
}

func getProperties(metadata *Metadata) map[string]string {
	properties := make(map[string]string)
	if metadata.Properties != nil && len(metadata.Properties.Properties) > 0 {
		properties = resolveProperties(metadata)
	}
	return properties
}

func createComponent(metadata *Metadata, path string) *cyclonedx.Component {
	c := component.New(metadata.ArtifactID, metadata.Version, Type)
	component.AddOrigin(c, path)
	component.AddType(c, Type)
	component.AddDescription(c, metadata.Description)
	return c
}

func addCPEs(c *cyclonedx.Component) {
	cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
	if len(cpes) > 0 {
		for _, cpe := range cpes {
			component.AddCPE(c, cpe)
		}
	}
}

func addComponentDetails(c *cyclonedx.Component, path, description string, metadata *Metadata) {
	component.AddOrigin(c, path)
	component.AddType(c, Type)
	component.AddDescription(c, description)

	meta := Dependency{
		GroupID:    metadata.GroupID,
		ArtifactID: metadata.ArtifactID,
		Version:    metadata.Version,
	}

	rawMetadata, err := helper.ToJSON(meta)
	if err != nil {
		log.Debugf(ErrConvertMetadataToJSON, err)
	}

	if len(rawMetadata) > 0 {
		component.AddRawMetadata(c, rawMetadata)
	}
}

func processDependencies(metadata *Metadata, properties map[string]string, path string, addr *urn.URN) {
	dependencyNode := &cyclonedx.Dependency{
		Ref:          metadata.ArtifactID,
		Dependencies: &[]string{},
	}

	for _, dependency := range metadata.Dependencies {
		resolveDependencyVersion(&dependency, metadata, properties)

		if dependency.ArtifactID == "" || dependency.Version == "" {
			continue
		}

		if strings.Contains(dependency.Version, "${") || strings.Contains(dependency.ArtifactID, "${") {
			continue
		}

		c := component.New(dependency.ArtifactID, dependency.Version, Type)
		addCPEs(c)
		c.PackageURL = MavenPackageURLPrefix + dependency.GroupID + "/" + dependency.ArtifactID + "@" + dependency.Version
		c.BOMRef = c.PackageURL

		component.AddOrigin(c, path)

		if len(dependency.GroupID) > 0 {
			c.Group = dependency.GroupID
		}

		rawMetadata, err := helper.ToJSON(dependency)
		if err != nil {
			log.Debugf(ErrConvertMetadataToJSON, err)
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

func resolveDependencyVersion(dependency *Dependency, metadata *Metadata, properties map[string]string) {
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
}

//nolint:all
// func readManifestFile(manifest types.ManifestFile, addr *urn.URN) {
// 	metadata, err := parseManifestFile(manifest.Content)
// 	if err != nil {
// 		log.Debugf("Failed to parse manifest file: %v", err)
// 		return
// 	}

// 	if len(metadata["Bundle-SymbolicName"]) == 0 || len(metadata["Bundle-Version"]) == 0 {
// 		return
// 	}

// 	c := component.New(metadata["Bundle-SymbolicName"], metadata["Bundle-Version"], Type)

// 	cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
// 	if len(cpes) > 0 {
// 		for _, cpe := range cpes {
// 			component.AddCPE(c, cpe)
// 		}
// 	}

// 	c.PackageURL = MavenPackageURLPrefix + c.Name + "@" + c.Version

// 	component.AddOrigin(c, manifest.Path)
// 	component.AddType(c, Type)

// 	rawMetadata, err := helper.ToJSON(metadata)
// 	if err != nil {
// 		log.Debugf(ErrConvertMetadataToJSON, err)
// 	}

// 	if len(rawMetadata) > 0 {
// 		component.AddRawMetadata(c, rawMetadata)
// 	}

// 	cdx.AddComponent(c, addr)

// }

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
	c.PackageURL = MavenPackageURLPrefix + c.Name + "@" + c.Version

	component.AddOrigin(c, manifest.Path)
	component.AddType(c, Type)

	rawMetadata, err := helper.ToJSON(metadata)
	if err != nil {
		log.Debugf(ErrConvertMetadataToJSON, err)
	}

	if len(rawMetadata) > 0 {
		component.AddRawMetadata(c, rawMetadata)
	}

	if len(layer) > 0 {
		component.AddLayer(c, layer)
	}

	cdx.AddComponent(c, addr)

}
