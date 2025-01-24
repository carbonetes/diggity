package npm

import (
	"strconv"
	"strings"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/types"
)

func scanPnpmLockfile(payload types.Payload) *[]cyclonedx.Component {
	components := []cyclonedx.Component{}
	manifest := payload.Body.(types.ManifestFile)
	metadata := readPnpmLockfile(manifest.Content)
	if metadata == nil {
		return nil
	}

	components = append(components, processPnpmDependencies(metadata.Dependencies, manifest, payload)...)
	components = append(components, processPackages(metadata.Packages, metadata.LockFileVersion, manifest, payload)...)

	return &components
}

func processPnpmDependencies(dependencies map[string]interface{}, manifest types.ManifestFile, payload types.Payload) []cyclonedx.Component {
	components := []cyclonedx.Component{}

	for name, info := range dependencies {
		if name == "" {
			continue
		}

		version := extractVersion(info)
		if len(version) == 0 {
			continue
		}

		c := createComponent(name, version, manifest, payload)
		components = append(components, *c)
	}

	return components
}

func processPackages(packages map[string]interface{}, lockFileVersion string, manifest types.ManifestFile, payload types.Payload) []cyclonedx.Component {
	components := []cyclonedx.Component{}
	separator := determineSeparator(lockFileVersion)

	for id := range packages {
		id = packageNameRegex.ReplaceAllString(id, "$1")
		id = strings.TrimPrefix(id, "/")
		props := strings.Split(id, separator)

		name := strings.Join(props[:len(props)-1], separator)
		version := props[len(props)-1]

		if len(name) == 0 || len(version) == 0 {
			continue
		}

		c := createComponent(name, version, manifest, payload)
		components = append(components, *c)
	}

	return components
}

func extractVersion(info interface{}) string {
	var version string
	switch t := info.(type) {
	case string:
		version = strings.SplitN(t, "()", 2)[0]
	case map[string]interface{}:
		val, ok := t["version"].(string)
		if ok {
			version = strings.SplitN(val, "(", 2)[0]
		}
	}
	return version
}

func createComponent(name, version string, manifest types.ManifestFile, payload types.Payload) *cyclonedx.Component {
	c := component.New(name, version, Type)

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

	return c
}

func determineSeparator(lockFileVersion string) string {
	separator := "/"
	version, _ := strconv.ParseFloat(lockFileVersion, 64)
	if version >= 6.0 {
		separator = "@"
	}
	return separator
}
