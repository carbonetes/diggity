package swift

import (
	"path/filepath"
	"slices"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "swift"

var Manifests = []string{"Package.resolved", ".package.resolved"}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	payload, ok := data.(types.Payload)
	if !ok {
		log.Debug("Swift Package Manager Handler received unknown type")
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

	metadata := readManifestFile(file.Content)

	switch metadata.Version {
	case 1:
		processPins(metadata.Object.Pins, file, payload)
	case 2:
		processPins(metadata.Pins, file, payload)
	}
}

func processPins(pins []Pin, file types.ManifestFile, payload types.Payload) {
	for _, pin := range pins {
		name, version := pin.Name, pin.State.Version

		c := component.New(name, version, Type)

		addCPEs(c)
		addComponentDetails(c, file, payload, pin)
		cdx.AddComponent(c, payload.Address)
	}
}

func addCPEs(c *cyclonedx.Component) {
	cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
	if len(cpes) > 0 {
		for _, cpe := range cpes {
			component.AddCPE(c, cpe)
		}
	}
}

func addComponentDetails(c *cyclonedx.Component, file types.ManifestFile, payload types.Payload, pin Pin) {
	component.AddOrigin(c, file.Path)
	component.AddType(c, Type)

	rawMetadata, err := helper.ToJSON(pin)
	if err != nil {
		log.Debugf("Error converting metadata to JSON: %s", err)
	}

	if len(rawMetadata) > 0 {
		component.AddRawMetadata(c, rawMetadata)
	}

	if len(payload.Layer) > 0 {
		component.AddLayer(c, payload.Layer)
	}
}
