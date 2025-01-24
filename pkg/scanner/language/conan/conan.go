package conan

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
	"github.com/golistic/urn"
)

const Type string = "conan"

var Manifests = []string{"conanfile.txt", "conan.lock"}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	payload, ok := data.(types.Payload)
	if !ok {
		log.Debug("Conan Handler received unknown type")
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

	if strings.Contains(file.Path, "conanfile.txt") {
		handleConanfileTxt(file, payload)
	} else if strings.Contains(file.Path, "conan.lock") {
		handleConanLock(file, payload)
	}
}

func handleConanfileTxt(file types.ManifestFile, payload types.Payload) {
	packages := readManifestFile(file.Content)
	if len(packages) == 0 {
		return
	}

	for _, pkg := range packages {
		if pkg.Name == "" {
			continue
		}

		c := component.New(pkg.Name, pkg.Version, Type)
		addCPEs(c, c.Name, c.Version)
		addCommonAttributes(c, file.Path, payload.Layer, pkg, payload.Address)
	}
}

func handleConanLock(file types.ManifestFile, payload types.Payload) {
	metadata := readLockFile(file.Content)
	if metadata == nil || len(metadata.GraphLock.Nodes) == 0 {
		return
	}

	for _, node := range metadata.GraphLock.Nodes {
		if node.Ref == "" {
			continue
		}

		name, version := parseNodeRef(node.Ref)
		if name == "" || version == "" {
			continue
		}

		c := component.New(name, version, Type)
		addCPEs(c, c.Name, c.Version)
		addCommonAttributes(c, file.Path, payload.Layer, node, payload.Address)
	}
}

func addCPEs(c *cyclonedx.Component, name, version string) {
	cpes := cpe.NewCPE23(name, name, version, Type)
	if len(cpes) > 0 {
		for _, cpe := range cpes {
			component.AddCPE(c, cpe)
		}
	}
}

func addCommonAttributes(c *cyclonedx.Component, filePath, layer string, metadata interface{}, address *urn.URN) {
	component.AddOrigin(c, filePath)
	component.AddType(c, Type)

	rawMetadata, err := helper.ToJSON(metadata)
	if err != nil {
		log.Debugf("Error converting metadata to JSON: %s", err)
	}

	if len(rawMetadata) > 0 {
		component.AddRawMetadata(c, rawMetadata)
	}

	if len(layer) > 0 {
		component.AddLayer(c, layer)
	}

	cdx.AddComponent(c, address)
}

func parseNodeRef(ref string) (string, string) {
	var val []string
	if strings.Contains(ref, "@") {
		val = strings.Split(ref, "@")
	} else if strings.Contains(ref, "#") {
		val = strings.Split(ref, "#")
	} else {
		val = []string{ref}
	}

	parts := strings.Split(val[0], "/")
	if len(parts) != 2 {
		return "", ""
	}

	name, version := parts[0], parts[1]
	if strings.Contains(version, "[") && strings.Contains(version, "]") {
		version = strings.Replace(version, "[", "", -1)
		version = strings.Replace(version, "]", "", -1)
	}

	return name, version
}
