package cargo

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

const Type string = "rust-crate"

var (
	Manifests = []string{"Cargo.lock"}
)

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	payload, ok := data.(types.Payload)
	if !ok {
		log.Debug("Cargo Handler received unknown type")
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

	packages := readManifestFile(file.Content)

	for _, pkg := range packages {
		processPackage(pkg, file.Path, payload)
	}
}

func processPackage(pkg, filePath string, payload types.Payload) {
	if !strings.Contains(pkg, "[[package]]") {
		return
	}

	metadata := parseMetadata(pkg)
	if metadata == nil {
		return
	}

	if metadata["Name"] == nil {
		return
	}

	c := component.New(metadata["Name"].(string), metadata["Version"].(string), Type)

	addCPEs(c, c.Name, c.Version)
	component.AddOrigin(c, filePath)
	component.AddType(c, Type)

	addRawMetadata(c, metadata)
	addLayer(c, payload.Layer)

	cdx.AddComponent(c, payload.Address)
}

func addCPEs(c *cyclonedx.Component, name, version string) {
	cpes := cpe.NewCPE23(name, name, version, Type)
	if len(cpes) > 0 {
		for _, cpe := range cpes {
			component.AddCPE(c, cpe)
		}
	}
}

func addRawMetadata(c *cyclonedx.Component, metadata map[string]interface{}) {
	rawMetadata, err := helper.ToJSON(metadata)
	if err != nil {
		log.Debugf("Error converting metadata to JSON: %s", err)
	}

	if len(rawMetadata) > 0 {
		component.AddRawMetadata(c, rawMetadata)
	}
}

func addLayer(c *cyclonedx.Component, layer string) {
	if len(layer) > 0 {
		component.AddLayer(c, layer)
	}
}
