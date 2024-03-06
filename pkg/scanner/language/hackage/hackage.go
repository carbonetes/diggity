package hackage

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
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
	manifest, ok := data.(types.ManifestFile)
	if !ok {
		log.Error("Hackage Handler received unknown type")
	}

	if strings.Contains(manifest.Path, "stack.yaml") {
		stackConfig := readStackConfigFile(manifest.Content)
		for _, dep := range stackConfig.ExtraDeps {
			name, version, _, _, _ := parseExtraDep(dep)
			if name == "" || version == "" {
				continue
			}

			c := component.New(name, version, Type)

			cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
			if len(cpes) > 0 {
				for _, cpe := range cpes {
					component.AddCPE(c, cpe)
				}
			}

			component.AddOrigin(c, manifest.Path)
			component.AddType(c, Type)

			cdx.AddComponent(c)
		}
	} else if strings.Contains(manifest.Path, "stack.yaml.lock") {
		lockFile := readStackLockConfigFile(manifest.Content)

		for _, pkg := range lockFile.Packages {
			name, version, _, _, _ := parseExtraDep(pkg.Original.Hackage)
			if name == "" || version == "" {
				continue
			}

			c := component.New(name, version, Type)

			cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
			if len(cpes) > 0 {
				for _, cpe := range cpes {
					component.AddCPE(c, cpe)
				}
			}

			component.AddOrigin(c, manifest.Path)
			component.AddType(c, Type)

			cdx.AddComponent(c)
		}
	} else if strings.Contains(manifest.Path, "cabal.project.freeze") {
		packages := readManifestFile(manifest.Content)
		for _, pkg := range packages {
			name, version, _, _, _ := parseExtraDep(pkg)
			if name == "" || version == "" {
				continue
			}

			c := component.New(name, version, Type)

			cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
			if len(cpes) > 0 {
				for _, cpe := range cpes {
					component.AddCPE(c, cpe)
				}
			}

			component.AddOrigin(c, manifest.Path)
			component.AddType(c, Type)

			cdx.AddComponent(c)
		}
	}

	return data
}

// Format Name and Version for parsing
func formatCabalPackage(anyPkg string) string {
	pkg := strings.Replace(strings.TrimSpace(anyPkg), "any.", "", -1)
	nv := strings.Replace(pkg, " ==", "-", -1)
	return strings.Replace(nv, ",", "", -1)
}
