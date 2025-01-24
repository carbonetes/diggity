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

var Manifests = []string{"buildscript-gradle.lockfile", ".build.gradle"}

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

	lines := strings.Split(string(file.Content), "\n")
	for _, line := range lines {
		processLine(line, file, payload)
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
