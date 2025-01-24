package npm

import (
	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/types"
)

func scanPackageJSON(payload types.Payload) *[]cyclonedx.Component {
	components := []cyclonedx.Component{}
	manifest := payload.Body.(types.ManifestFile)
	metadata := readManifestFile(manifest.Content)
	if metadata == nil {
		return nil
	}

	processDependencies(metadata.DevDependencies, manifest, payload, &components)
	processDependencies(metadata.Dependencies, manifest, payload, &components)

	if metadata.Name != "" && metadata.Version != "" {
		processMainComponent(metadata, manifest, payload, &components)
	}

	return &components
}

func processDependencies(dependencies map[string]interface{}, manifest types.ManifestFile, payload types.Payload, components *[]cyclonedx.Component) {
	for name, version := range dependencies {
		n := parseYarnPackageName(name)
		v := cleanVersion(version.(string))
		if n == "" || v == "" || !validateVersion(v) {
			continue
		}

		c := component.New(n, v, Type)
		addComponentDetails(c, manifest, payload)
		*components = append(*components, *c)
	}
}

func processMainComponent(metadata *Metadata, manifest types.ManifestFile, payload types.Payload, components *[]cyclonedx.Component) {
	n := cleanName(metadata.Name)
	v := cleanVersion(metadata.Version)

	c := component.New(n, v, Type)
	addComponentDetails(c, manifest, payload)

	switch metadata.License.(type) {
	case string:
		component.AddLicense(c, metadata.License.(string))
	case map[string]interface{}:
		license := metadata.License.(map[string]interface{})
		if _, ok := license["type"]; ok {
			component.AddLicense(c, license["type"].(string))
		}
	}

	rawMetadata, err := helper.ToJSON(metadata)
	if err != nil {
		log.Debugf("Error converting metadata to JSON: %s", err)
	}

	if len(rawMetadata) > 0 {
		component.AddRawMetadata(c, rawMetadata)
	}

	*components = append(*components, *c)
}

func addComponentDetails(c *cyclonedx.Component, manifest types.ManifestFile, payload types.Payload) {
	cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
	for _, cpe := range cpes {
		component.AddCPE(c, cpe)
	}

	component.AddOrigin(c, manifest.Path)
	component.AddType(c, Type)

	if len(payload.Layer) > 0 {
		component.AddLayer(c, payload.Layer)
	}
}

func scanPackageLockfile(payload types.Payload) *[]cyclonedx.Component {
	components := []cyclonedx.Component{}
	manifest := payload.Body.(types.ManifestFile)
	metadata := readPackageLockfile(manifest.Content)
	if metadata == nil {
		return nil
	}

	if len(metadata.Dependencies) == 0 {
		return nil
	}
	for name, dependency := range metadata.Dependencies {
		c := component.New(parseYarnPackageName(name), cleanVersion(dependency.Version), Type)

		if len(c.Name) == 0 || len(c.Version) == 0 {
			continue
		}

		cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
		if len(cpes) > 0 {
			for _, cpe := range cpes {
				component.AddCPE(c, cpe)
			}
		}

		component.AddOrigin(c, manifest.Path)
		component.AddType(c, Type)

		if len(payload.Layer) > 0 {
			component.AddLayer(c, payload.Layer)
		}

		components = append(components, *c)
	}

	return &components
}
