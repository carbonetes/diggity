package golang

import (
	"bytes"
	"debug/buildinfo"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/types"
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

	for _, s := range buildInfo.Settings {
		// locate version
		v := parseVersion(s.Value)
		if v != "" {
			c := component.New(file.Path, v, Type)

			cpes := GenerateCpes(c.Version, SplitPath(c.Name))
			if len(cpes) > 0 {
				for _, cpe := range cpes {
					component.AddCPE(c, cpe)
				}
			}

			component.AddOrigin(c, file.Path)
			component.AddType(c, Type)

			rawMetadata, err := helper.ToJSON(s)
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
			break
		}
	}

	for _, dep := range buildInfo.Deps {
		if dep.Path == "" || dep.Version == "" {
			continue
		}

		c := component.New(dep.Path, dep.Version, Type)

		cpes := GenerateCpes(c.Version, SplitPath(c.Name))
		if len(cpes) > 0 {
			for _, cpe := range cpes {
				component.AddCPE(c, cpe)
			}
		}

		component.AddOrigin(c, file.Path)
		component.AddType(c, Type)

		rawMetadata, err := helper.ToJSON(dep)
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
