package nuget

import (
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/types"

	"github.com/saferwall/pe"
)

func parsePE(data []byte) (*pe.File, bool) {
	file, err := pe.NewBytes(data, &pe.Options{})
	if err != nil {
		return nil, false
	}

	err = file.Parse()
	if err != nil {
		return nil, false
	}

	return file, true
}

func scanPE(payload types.Payload, peFile *pe.File) {
	// parse version resource from the PE file
	versionInfo, err := peFile.ParseVersionResources()
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
	component.AddOrigin(c, payload.Body.(types.ManifestFile).Path)
	component.AddType(c, Type)

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
