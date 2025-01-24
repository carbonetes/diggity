package nuget

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

const Type string = "nuget"

var Manifests = []string{".deps.json"}

var projectFiles = []string{".csproj", ".vbproj"}	

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}

	if filepath.Ext(file) == ".dll" || filepath.Ext(file) == ".exe" {
		return Type, true, true
	}

	if slices.Contains(projectFiles, filepath.Ext(file)) {
		return Type, true, true
	}

	return "", false, false
}

func Scan(data interface{}) interface{} {
	payload, ok := data.(types.Payload)
	if !ok {
		log.Debug("Nuget Handler received unknown type")
		return nil
	}

	if filepath.Ext(payload.Body.(types.ManifestFile).Path) == ".dll" || filepath.Ext(payload.Body.(types.ManifestFile).Path) == ".exe" {
		if peFile, isPE := parsePE(payload.Body.(types.ManifestFile).Content); isPE {
			scanPE(payload, peFile)
		}
		return data
	}

	if slices.Contains(projectFiles, filepath.Ext(payload.Body.(types.ManifestFile).Path)) {
		scanProjectFile(payload)
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

	metadata := readManifestFile(file.Content)

	if len(metadata.Libraries) == 0 {
		return
	}

	for id, pkg := range metadata.Libraries {
		processLibrary(id, pkg, file, payload)
	}
}

func processLibrary(id string, pkg DotnetLibrary, file types.ManifestFile, payload types.Payload) {
	if pkg.Type != "package" {
		return
	}

	attributes := strings.Split(id, "/")
	if len(attributes) != 2 {
		return
	}

	name, version := attributes[0], attributes[1]

	if name == "" || version == "" {
		return
	}

	c := component.New(name, version, Type)

	addCPEs(c, name, version)
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

func addCPEs(c *cyclonedx.Component, name, version string) {
	cpes := cpe.NewCPE23(name, name, version, Type)
	if len(cpes) > 0 {
		for _, cpe := range cpes {
			component.AddCPE(c, cpe)
		}
	}
}
