package hackage

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

const Type string = "hackage"

var Manifests = []string{"cabal.project.freeze", "stack.yaml", "stack.yaml.lock"}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	payload, ok := data.(types.Payload)
	if !ok {
		log.Debug("Hackage Handler received unknown type")
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

	switch {
	case strings.Contains(file.Path, "stack.yaml"):
		handleStackYaml(file, payload)
	case strings.Contains(file.Path, "stack.yaml.lock"):
		handleStackYamlLock(file, payload)
	case strings.Contains(file.Path, "cabal.project.freeze"):
		handleCabalProjectFreeze(file, payload)
	}
}

func handleStackYaml(file types.ManifestFile, payload types.Payload) {
	stackConfig := readStackConfigFile(file.Content)
	if stackConfig == nil {
		return
	}

	for _, dep := range stackConfig.ExtraDeps {
		name, version, _, _, _ := parseExtraDep(dep)
		if name == "" || version == "" {
			continue
		}

		c := component.New(name, version, Type)
		addComponentDetails(c, dep, file.Path, payload)
	}
}

func handleStackYamlLock(file types.ManifestFile, payload types.Payload) {
	lockFile := readStackLockConfigFile(file.Content)
	if lockFile == nil {
		return
	}

	for _, pkg := range lockFile.Packages {
		name, version, _, _, _ := parseExtraDep(pkg.Original.Hackage)
		if name == "" || version == "" {
			continue
		}

		c := component.New(name, version, Type)
		addComponentDetails(c, pkg, file.Path, payload)
	}
}

func handleCabalProjectFreeze(file types.ManifestFile, payload types.Payload) {
	packages := readManifestFile(file.Content)
	if len(packages) == 0 {
		return
	}

	for _, pkg := range packages {
		name, version, _, _, _ := parseExtraDep(pkg)
		if name == "" || version == "" {
			continue
		}

		c := component.New(name, version, Type)
		addComponentDetails(c, pkg, file.Path, payload)
	}
}

func addComponentDetails(c *cyclonedx.Component, metadata interface{}, filePath string, payload types.Payload) {
	cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
	if len(cpes) > 0 {
		for _, cpe := range cpes {
			component.AddCPE(c, cpe)
		}
	}

	rawMetadata, err := helper.ToJSON(metadata)
	if err != nil {
		log.Debugf("Error converting metadata to JSON: %s", err)
	}

	component.AddOrigin(c, filePath)
	component.AddType(c, Type)
	component.AddRawMetadata(c, rawMetadata)

	if len(payload.Layer) > 0 {
		component.AddLayer(c, payload.Layer)
	}

	cdx.AddComponent(c, payload.Address)
}

// Format Name and Version for parsing
func formatCabalPackage(anyPkg string) string {
	pkg := strings.Replace(strings.TrimSpace(anyPkg), "any.", "", -1)
	nv := strings.Replace(pkg, " ==", "-", -1)
	return strings.Replace(nv, ",", "", -1)
}
