package hex

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

const Type string = "hex"

var Manifests = []string{"rebar.lock", "mix.lock"}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	payload, ok := data.(types.Payload)
	if !ok {
		log.Debug("Hex Handler received unknown type")
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

	if strings.Contains(file.Path, "rebar.lock") {
		processRebarFile(file, payload)
	} else if strings.Contains(file.Path, "mix.lock") {
		processMixFile(file, payload)
	}
}

func processRebarFile(file types.ManifestFile, payload types.Payload) {
	metadata := readRebarFile(file.Content)
	if len(metadata) == 0 {
		return
	}

	for _, pkg := range metadata {
		if pkg.Name == "" || pkg.Version == "" {
			continue
		}

		c := createComponent(pkg, file.Path)
		addCPEs(c, c.Name, c.Version)
		addMetadata(c, pkg, payload.Layer)
		cdx.AddComponent(c, payload.Address)
	}
}

func processMixFile(file types.ManifestFile, payload types.Payload) {
	metadata := readMixFile(file.Content)
	if len(metadata) == 0 {
		return
	}

	for _, m := range metadata {
		if m.Name == "" || m.Version == "" {
			continue
		}

		c := createComponent(m, file.Path)
		addCPEs(c, c.Name, c.Version)
		addMetadata(c, m, payload.Layer)
		cdx.AddComponent(c, payload.Address)
	}
}

func createComponent(m HexMetadata, filePath string) *cyclonedx.Component {
	c := component.New(m.Name, m.Version, Type)
	component.AddOrigin(c, filePath)
	component.AddType(c, Type)
	return c
}

func addCPEs(c *cyclonedx.Component, name, version string) {
	cpes := cpe.NewCPE23(name, name, version, Type)
	if len(cpes) > 0 {
		for _, cpe := range cpes {
			component.AddCPE(c, cpe)
		}
	}
}

func addMetadata(c *cyclonedx.Component, m HexMetadata, layer string) {
	rawMetadata, err := helper.ToJSON(m)
	if err != nil {
		log.Debugf("Error converting metadata to JSON: %s", err)
	}

	if len(rawMetadata) > 0 {
		component.AddRawMetadata(c, rawMetadata)
	}

	if len(layer) > 0 {
		component.AddLayer(c, layer)
	}
}
