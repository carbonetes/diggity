package hex

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
		packages := readRebarFile(file.Content)
		if len(packages) == 0 {
			return
		}

		for _, pkg := range packages {
			if pkg.Name == "" || pkg.Version == "" {
				continue
			}

			c := component.New(pkg.Name, pkg.Version, Type)

			cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
			if len(cpes) > 0 {
				for _, cpe := range cpes {
					component.AddCPE(c, cpe)
				}
			}

			component.AddOrigin(c, file.Path)
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

	} else if strings.Contains(file.Path, "mix.lock") {
		packages := readMixFile(file.Content)
		if len(packages) == 0 {
			return
		}

		for _, pkg := range packages {
			if pkg.Name == "" || pkg.Version == "" {
				continue
			}

			c := component.New(pkg.Name, pkg.Version, Type)

			cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
			if len(cpes) > 0 {
				for _, cpe := range cpes {
					component.AddCPE(c, cpe)
				}
			}

			component.AddOrigin(c, file.Path)
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
	}
}
