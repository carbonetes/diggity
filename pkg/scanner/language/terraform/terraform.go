package terraform

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

const Type string = "terraform"

var Manifests = []string{".terraform.lock.hcl", "terraform.lock.hcl"}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	payload, ok := data.(types.Payload)
	if !ok {
		log.Debug("Terraform Handler received unknown type")
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

	metadata, err := readLockfile(file.Content)
	if err != nil {
		log.Debugf("Error reading lockfile: %s", err.Error())
		return
	}

	if metadata == nil || len(metadata.Providers) == 0 {
		return
	}

	processProviders(metadata.Providers, file, payload)
}

func processProviders(providers []Provider, file types.ManifestFile, payload types.Payload) {
	for _, provider := range providers {
		name, version := provider.URL, provider.Version
		if name == "" || version == "" {
			continue
		}

		c := createComponent(name, version, file, payload)
		cdx.AddComponent(c, payload.Address)
	}
}

func createComponent(name, version string, file types.ManifestFile, payload types.Payload) *cyclonedx.Component {
	c := component.New(name, version, Type)

	addCPEs(c)
	component.AddOrigin(c, file.Path)
	component.AddType(c, Type)
	addRawMetadata(c, file)
	addLayer(c, payload)

	return c
}

func addCPEs(c *cyclonedx.Component) {
	cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
	for _, cpe := range cpes {
		component.AddCPE(c, cpe)
	}
}

func addRawMetadata(c *cyclonedx.Component, file types.ManifestFile) {
	rawMetadata, err := helper.ToJSON(file)
	if err != nil {
		log.Debugf("Error converting metadata to JSON: %s", err)
		return
	}
	if len(rawMetadata) > 0 {
		component.AddRawMetadata(c, rawMetadata)
	}
}

func addLayer(c *cyclonedx.Component, payload types.Payload) {
	if len(payload.Layer) > 0 {
		component.AddLayer(c, payload.Layer)
	}
}
