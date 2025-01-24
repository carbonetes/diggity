package cocoapods

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

const Type string = "cocoapods"

var Manifests = []string{"Podfile.lock"}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	payload, ok := data.(types.Payload)
	if !ok {
		log.Debug("Cocoapods Handler received unknown type")
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
	for _, pod := range metadata.Pods {
		processPod(pod, file, metadata, payload)
	}
}

func processPod(pod interface{}, file types.ManifestFile, metadata FileLockMetadata, payload types.Payload) {
	var pods string
	switch c := pod.(type) {
	case string:
		pods = c
	case map[string]interface{}:
		val := pod.(map[string]interface{})
		for all := range val {
			pods = all
		}
	}

	attributes := strings.Split(pods, " ")
	name, version := attributes[0], strings.TrimSuffix(strings.TrimPrefix(attributes[1], "("), ")")

	c := component.New(name, version, Type)

	addCPEs(c, Type)
	addMetadata(c, file, metadata, payload)
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

func addMetadata(c *cyclonedx.Component, file types.ManifestFile, metadata FileLockMetadata, payload types.Payload) {
	rawMetadata, err := helper.ToJSON(metadata)
	if err != nil {
		log.Debug("Failed to convert metadata to JSON")
	}

	component.AddOrigin(c, file.Path)
	component.AddType(c, Type)

	if len(rawMetadata) > 0 {
		component.AddRawMetadata(c, rawMetadata)
	}

	if len(payload.Layer) > 0 {
		component.AddLayer(c, payload.Layer)
	}
}
