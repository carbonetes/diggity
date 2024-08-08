package golang

import (
	"path/filepath"
	"slices"

	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "golang"

var Manifests = []string{"go.mod"}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}

	if len(filepath.Ext(file)) == 0 {
		return Type, true, true
	}

	return "", false, false
}

func Scan(data interface{}) interface{} {
	payload, ok := data.(types.Payload)
	if !ok {
		log.Debug("Go Modules Handler received unknown type")
		return nil
	}

	if len(filepath.Ext(payload.Body.(types.ManifestFile).Path)) == 0 {
		if bInfo, isGoBin := parseGoBin(payload.Body.(types.ManifestFile).Content); isGoBin {
			scanBinary(payload, bInfo)
		}
		return data
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

	modFile := readManifestFile(file.Content, file.Path)
	if modFile == nil {
		return
	}

	for _, pkg := range modFile.Require {
		if pkg.Mod.Path == "" || pkg.Mod.Version == "" {
			continue
		}
		if checkIfExcluded(modFile.Exclude, pkg.Mod.Path) {
			continue
		}

		c := component.New(pkg.Mod.Path, pkg.Mod.Version, Type)

		cpes := GenerateCpes(c.Version, SplitPath(c.Name))
		if len(cpes) > 0 {
			for _, cpe := range cpes {
				component.AddCPE(c, cpe)
			}
		}

		rawMetadata, err := helper.ToJSON(pkg)
		if err != nil {
			log.Debugf("Error converting metadata to JSON: %s", err)
		}

		component.AddOrigin(c, file.Path)
		component.AddType(c, Type)
		component.AddRawMetadata(c, rawMetadata)

		cdx.AddComponent(c, payload.Address)
	}

	for _, pkg := range modFile.Replace {
		if pkg.New.Path == "" || pkg.New.Version == "" {
			continue
		}
		if checkIfExcluded(modFile.Exclude, pkg.New.Path) {
			continue
		}

		c := component.New(pkg.New.Path, pkg.New.Version, Type)

		cpes := GenerateCpes(c.Version, SplitPath(c.Name))
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
