package pacman

import (
	"path/filepath"
	"slices"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "archlinux"

var Manifests = []string{"var/lib/pacman/local"}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, file) && filepath.Base(file) == "desc" {
		return Type, true, true
	}
	return "", false, true
}

func Scan(data interface{}) interface{} {
	payload, ok := data.(types.Payload)
	if !ok {
		log.Debug("Pacman Handler received unknown type")
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

	attributes := helper.SplitContentsByEmptyLine(string(file.Content))
	metadata := parseMetadata(attributes)

	if metadata["name"] == nil || metadata["name"] == "" {
		return
	}

	name, version, desc := metadata["name"].(string), metadata["version"].(string), metadata["description"].(string)

	c := component.New(name, version, Type)

	cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
	if len(cpes) > 0 {
		for _, cpe := range cpes {
			component.AddCPE(c, cpe)
		}
	}

	component.AddOrigin(c, file.Path)
	component.AddType(c, Type)
	component.AddDescription(c, desc)

	arch, ok := metadata["arch"].(string)
	if !ok {
		arch = ""
	}

	c.PackageURL = c.PackageURL + "?arch=" + arch

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
