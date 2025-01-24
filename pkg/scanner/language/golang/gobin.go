package golang

import (
	"bytes"
	"debug/buildinfo"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/types"
	"github.com/golistic/urn"
	"golang.org/x/mod/modfile"
)

func parseGoBin(data []byte) (*buildinfo.BuildInfo, bool) {
	build, err := buildinfo.Read(bytes.NewReader(data))
	if err != nil {
		return nil, false
	}

	return build, true
}

func scanBinary(payload types.Payload, buildInfo *debug.BuildInfo) {
	file, ok := payload.Body.(types.ManifestFile)
	if !ok {
		log.Debugf("Failed to convert payload body to manifest file")
		return
	}

	processSettings(file, payload, buildInfo.Settings)
	processDependencies(file, payload, buildInfo.Deps)
}

func processSettings(file types.ManifestFile, payload types.Payload, settings []debug.BuildSetting) {
	for _, s := range settings {
		// locate version
		v := parseVersion(s.Value)
		if v != "" {
			c := component.New(file.Path, v, Type)

			addCPEs(c)
			addComponentDetails(c, file.Path, payload.Layer, s, payload.Address)
			break
		}
	}
}

func processDependencies(file types.ManifestFile, payload types.Payload, deps []*debug.Module) {
	for _, dep := range deps {
		if dep.Path == "" || dep.Version == "" {
			continue
		}

		c := component.New(dep.Path, dep.Version, Type)

		addCPEs(c)
		addComponentDetails(c, file.Path, payload.Layer, dep, payload.Address)
	}
}

func addCPEs(c *cyclonedx.Component) {
	cpes := GenerateCpes(c.Version, SplitPath(c.Name))
	if len(cpes) > 0 {
		for _, cpe := range cpes {
			component.AddCPE(c, cpe)
		}
	}
}

func addComponentDetails(c *cyclonedx.Component, filePath, layer string, metadata interface{}, address *urn.URN) {
	component.AddOrigin(c, filePath)
	component.AddType(c, Type)

	rawMetadata, err := helper.ToJSON(metadata)
	if err != nil {
		log.Debugf("Error converting metadata to JSON: %s", err)
	}

	if len(rawMetadata) > 0 {
		component.AddRawMetadata(c, rawMetadata)
	}

	if len(layer) > 0 {
		component.AddLayer(c, layer)
	}

	cdx.AddComponent(c, address)
}

func checkIfExcluded(excludeList []*modfile.Exclude, name string) bool {
	for _, exclude := range excludeList {
		if exclude.Mod.Path == name {
			return true
		}
	}
	return false
}

// Split go package path
func SplitPath(path string) []string {
	return strings.Split(path, "/")
}

func GenerateCpes(version string, paths []string) []string {
	var cpes []string
	// check if cpePaths only contains the product
	if len(paths) > 1 {
		cpes = cpe.NewCPE23(paths[len(paths)-2], paths[len(paths)-1], FormatVersion(version), Type)
	} else {
		cpes = cpe.NewCPE23("", paths[0], FormatVersion(version), Type)
	}
	return cpes
}

// Format Version String
func FormatVersion(version string) string {
	if strings.Contains(version, "(") && strings.Contains(version, ")") {
		version = strings.Replace(version, "(", "", -1)
		version = strings.Replace(version, ")", "", -1)
	}
	return version
}

var pattern = regexp.MustCompile(`\d+\.\d+\.\d+`)

func parseVersion(s string) string {
	return pattern.FindString(s)
}
