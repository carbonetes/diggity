package nuget

import (
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/scanner/binary/pe"
	"github.com/carbonetes/diggity/pkg/types"
)

func scanPE(payload types.Payload) {
	peFile := payload.Body.(*pe.PEFile)

	// parse version resource from the PE file
	versionInfo, err := peFile.File.ParseVersionResources()
	if err != nil {
		log.Debug(err)
		return
	}

	var name, version string
	if versionInfo != nil {
		if v, ok := versionInfo["FileDescription"]; ok {
			name = v
		}

		if v, ok := versionInfo["FileVersion"]; ok {
			version = v
		}
	}

	if name == "" || version == "" {
		return
	}

	// create a new component
	c := component.New(name, version, Type)

	component.AddLayer(c, payload.Layer)
	component.AddOrigin(c, peFile.Path)

	rawMetadata, err := helper.ToJSON(versionInfo)
	if err != nil {
		log.Debugf("Error converting metadata to JSON: %s", err)
	}

	if len(rawMetadata) > 0 {
		component.AddRawMetadata(c, rawMetadata)
	}

	// add the component to the bom
	cdx.AddComponent(c, payload.Address)

}
