package cargo

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
		if !strings.Contains(pkg, "[[package]]") {
			continue
		}

		metadata := parseMetadata(pkg)
		if metadata == nil {
			continue
		}

		if metadata["Name"] == nil {
			return
		}

		c := component.New(metadata["Name"].(string), metadata["Version"].(string), Type)

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

		if len(payload.Layer) > 0 {
			component.AddLayer(c, payload.Layer)
		}

		cdx.AddComponent(c, payload.Address)
	}
}
