package gradle

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "gradle"

var Manifests = []string{"buildscript-gradle.lockfile", ".build.gradle", "libs.versions.toml"}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	payload, ok := data.(types.Payload)
	if !ok {
		log.Debug("Gradle Handler received unknown type")
		return nil
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

	switch {
	case strings.Contains(file.Path, "buildscript-gradle.lockfile"), strings.Contains(file.Path, ".build.gradle"):
		lines := strings.Split(string(file.Content), "\n")
		for _, line := range lines {
			processLine(line, file, payload)
		}
	case strings.Contains(file.Path, "libs.versions.toml"):
		metadata := read(string(file.Content))
		if metadata == nil {
			return
		}
		processVersionMetadata(metadata, file, payload)
	}

}

func processLine(line string, file types.ManifestFile, payload types.Payload) {
	if !strings.Contains(line, ":") {
		return
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	attributes := strings.SplitN(line, ":", 3)
	if len(attributes) < 3 {
		return
	}

	metadata := Metadata{
		Vendor:  attributes[0],
		Name:    attributes[1],
		Version: strings.ReplaceAll(attributes[2], "=classpath", ""),
	}

	c := component.New(metadata.Name, metadata.Version, Type)

	cpes := cpe.NewCPE23(metadata.Vendor, c.Name, c.Version, Type)
	if len(cpes) > 0 {
		for _, cpe := range cpes {
			component.AddCPE(c, cpe)
		}
	}

	component.AddOrigin(c, file.Path)
	component.AddType(c, Type)

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

	cdx.AddComponent(c, payload.Address)
}

func processVersionMetadata(metadata *VersionMetadata, file types.ManifestFile, payload types.Payload) {
	processLibraries(metadata, file, payload)
	processPlugins(metadata, file, payload)
}

func processLibraries(metadata *VersionMetadata, file types.ManifestFile, payload types.Payload) {
	for id, lib := range metadata.Libraries {
		name, version := extractLibraryDetails(id, lib, metadata.Versions)
		if name == "" || version == "" {
			continue
		}
		addComponent(name, version, lib, file, payload)
	}
}

func processPlugins(metadata *VersionMetadata, file types.ManifestFile, payload types.Payload) {
	for id, plugin := range metadata.Plugins {
		name, version := extractPluginDetails(id, plugin, metadata.Versions)
		if name == "" || version == "" {
			continue
		}
		addComponent(name, version, plugin, file, payload)
	}
}

func extractLibraryDetails(id string, lib Library, versions map[string]string) (string, string) {
	name := id
	if lib.Group != nil {
		name = *lib.Group
	}
	if lib.Module != nil {
		name = *lib.Module
	}
	version := resolveVersion(lib.Version, versions)
	return name, version
}

func extractPluginDetails(id string, plugin Plugin, versions map[string]string) (string, string) {
	name := id
	if plugin.ID != "" {
		name = plugin.ID
	}
	version := resolveVersion(plugin.Version, versions)
	return name, version
}

func resolveVersion(versionData interface{}, versions map[string]string) string {
	var ref, version string
	if v, ok := versionData.(string); ok {
		version = v
	}
	if v, ok := versionData.(map[string]interface{}); ok {
		if vf, ok := v["ref"].(string); ok {
			ref = vf
		}
	}
	if vv, ok := versions[ref]; ok && vv != "" {
		version = vv
	}
	return version
}

func addComponent(name, version string, metadata interface{}, file types.ManifestFile, payload types.Payload) {
	c := component.New(name, version, Type)
	if len(payload.Layer) > 0 {
		component.AddLayer(c, payload.Layer)
	}
	rawMetadata, err := helper.ToJSON(metadata)
	if err != nil {
		log.Debugf("Error converting metadata to JSON: %s", err)
	}
	if len(rawMetadata) > 0 {
		component.AddRawMetadata(c, rawMetadata)
	}
	component.AddOrigin(c, file.Path)
	cdx.AddComponent(c, payload.Address)
}
