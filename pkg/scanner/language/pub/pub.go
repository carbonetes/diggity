package pub

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

const Type string = "pub"

var Manifests = []string{"pubspec.yaml", "pubspec.lock"}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	payload, ok := data.(types.Payload)
	if !ok {
		log.Debug("Pub Handler received unknown type")
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

	if strings.Contains(file.Path, "pubspec.yaml") {
		processPubspecYaml(file, payload)
	} else if strings.Contains(file.Path, "pubspec.lock") {
		processPubspecLock(file, payload)
	}
}

func processPubspecYaml(file types.ManifestFile, payload types.Payload) {
	metadata := readManifestFile(file.Content)
	if metadata == nil {
		return
	}

	var name, version, license string

	name, ok := metadata["name"].(string)
	if !ok {
		return
	}
	version, ok = metadata["version"].(string)
	if !ok {
		version = "0.0.0"
	}

	if val, ok := metadata["license"].(string); ok {
		license = val
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

	rawMetadata, err := helper.ToJSON(metadata)
	if err != nil {
		log.Debugf("Error converting metadata to JSON: %s", err)
	}

	if len(rawMetadata) > 0 {
		component.AddRawMetadata(c, rawMetadata)
	}

	if license != "" {
		component.AddLicense(c, license)
	}

	if len(payload.Layer) > 0 {
		component.AddLayer(c, payload.Layer)
	}

	cdx.AddComponent(c, payload.Address)
}

func processPubspecLock(file types.ManifestFile, payload types.Payload) {
	metadata := readLockFile(file.Content)
	if metadata == nil || len(metadata.Packages) == 0 {
		return
	}

	for _, pkg := range metadata.Packages {
		if pkg.Description.Name == "" || pkg.Version == "" {
			continue
		}
		processPackage(pkg, file.Path, payload)
	}
}

func processPackage(pkg Package, filePath string, payload types.Payload) {
	c := component.New(pkg.Description.Name, pkg.Version, Type)

	cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
	if len(cpes) > 0 {
		for _, cpe := range cpes {
			component.AddCPE(c, cpe)
		}
	}

	component.AddOrigin(c, filePath)
	component.AddType(c, Type)

	rawMetadata, err := helper.ToJSON(pkg)
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
