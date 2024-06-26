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

	devDependencies := metadata.DevDependencies
	if len(devDependencies) > 0 {
		for name, version := range devDependencies {
			n := parseYarnPackageName(name)
			v := cleanVersion(version.(string))
			if n == "" || v == "" {
				continue
			}

			if !validateVersion(v) {
				continue
			}

			c := component.New(n, v, Type)

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
			// cdx.AddComponent(c, payload.Address)
		}
	}
	dependencies := metadata.Dependencies
	if len(dependencies) > 0 {
		for name, version := range dependencies {
			n := parseYarnPackageName(name)
			v := cleanVersion(version.(string))
			if n == "" || v == "" {
				continue
			}

			if !validateVersion(v) {
				continue
			}

			c := component.New(n, v, Type)

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
			// cdx.AddComponent(c, payload.Address)
		}
	}

	if metadata.Name == "" || metadata.Version == "" {
		return &components
	}

	n := cleanName(metadata.Name)
	v := cleanVersion(metadata.Version)

	c := component.New(n, v, Type)

	cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
	if len(cpes) > 0 {
		for _, cpe := range cpes {
			component.AddCPE(c, cpe)
		}
	}

	component.AddOrigin(c, manifest.Path)
	component.AddType(c, Type)

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

	if len(payload.Layer) > 0 {
		component.AddLayer(c, payload.Layer)
	}

	// cdx.AddComponent(c, payload.Address)
	components = append(components, *c)
	return &components
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
		if name == "" || dependency.Version == "" {
			continue
		}

		c := component.New(parseYarnPackageName(name), cleanVersion(dependency.Version), Type)

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

		// cdx.AddComponent(c, payload.Address)
		components = append(components, *c)
	}

	return &components
}
