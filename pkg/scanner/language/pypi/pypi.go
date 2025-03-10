package pypi

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
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "pypi"

var (
	Manifests  = []string{"METADATA", "requirements.txt", "poetry.lock", "PKG-INFO"}
	Extensions = []string{".egg-info"}
)

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) || slices.Contains(Extensions, filepath.Ext(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	payload, ok := data.(types.Payload)
	if !ok {
		log.Debug("Python Handler received unknown type")
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

	switch filepath.Base(file.Path) {
	case "METADATA", "PKG-INFO":
		processMetadataFile(file, payload)
	case "requirements.txt":
		processRequirementsFile(file, payload)
	case "poetry.lock":
		processPoetryLockFile(file, payload)
	default:
		if filepath.Ext(file.Path) == ".egg-info" {
			processMetadataFile(file, payload)
		}
	}
}

func processMetadataFile(file types.ManifestFile, payload types.Payload) {
	metadata := readManifestFile(file.Content)

	name, ok := metadata["Name"].(string)
	if !ok {
		return
	}

	version, ok := metadata["Version"].(string)
	if !ok {
		return
	}

	if name == "" || version == "" || strings.Contains(name, "=") || strings.Contains(version, "=") {
		return
	}

	c := component.New(name, version, Type)

	cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
	if len(cpes) > 0 {
		for _, cpe := range cpes {
			component.AddCPE(c, cpe)
		}
	}

	component.AddOrigin(c, file.Path)
	component.AddType(c, Type)

	if val, ok := metadata["Summary"].(string); ok {
		component.AddDescription(c, val)
	}

	if val, ok := metadata["License"].(string); ok {
		component.AddLicense(c, val)
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

	cdx.AddComponent(c, payload.Address)
}

func processRequirementsFile(file types.ManifestFile, payload types.Payload) {
	attributes := readRequirementsFile(file.Content)
	for _, attribute := range attributes {
		if len(attribute) != 2 {
			continue
		}

		name, version := attribute[0], attribute[1]

		if name == "" || version == "" || strings.Contains(name, "=") || strings.Contains(version, "=") {
			continue
		}

		c := component.New(name, version, Type)

		cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
		if len(cpes) > 0 {
			for _, cpe := range cpes {
				component.AddCPE(c, cpe)
			}
		}

		component.AddOrigin(c, file.Path)
		component.AddType(c, Type)

		if len(payload.Layer) > 0 {
			component.AddLayer(c, payload.Layer)
		}

		cdx.AddComponent(c, payload.Address)
	}
}

func processPoetryLockFile(file types.ManifestFile, payload types.Payload) {
	metadata := readPoetryLockFile(file.Content)
	if metadata == nil {
		return
	}

	for _, packageInfo := range metadata.Packages {
		processPackageInfo(packageInfo, file, payload)
	}
}

func processPackageInfo(packageInfo Package, file types.ManifestFile, payload types.Payload) {
	name, version := packageInfo.Name, packageInfo.Version

	if name == "" || version == "" || strings.Contains(name, "=") || strings.Contains(version, "=") {
		return
	}

	c := component.New(name, version, Type)

	addCPEs(c, Type)
	component.AddOrigin(c, file.Path)
	component.AddType(c, Type)

	rawMetadata, err := helper.ToJSON(packageInfo)
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

func addCPEs(c *cyclonedx.Component, Type string) {
	cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
	if len(cpes) > 0 {
		for _, cpe := range cpes {
			component.AddCPE(c, cpe)
		}
	}
}
