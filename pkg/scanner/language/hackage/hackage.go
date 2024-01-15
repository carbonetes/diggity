package hackage

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/stream"
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
			name, version, pkgHash, size, rev := parseExtraDep(dep)
			if name == "" || version == "" {
				continue
			}

			metadata := HackageMetadata{
				Name:     name,
				Version:  version,
				PkgHash:  pkgHash,
				Size:     size,
				Revision: rev,
			}

			component := types.NewComponent(name, version, Type, manifest.Path, "", metadata)
			cpes := cpe.NewCPE23(component.Name, component.Name, component.Version, Type)
			if len(cpes) > 0 {
				component.CPEs = append(component.CPEs, cpes...)
			}
			stream.AddComponent(component)
		}
	} else if strings.Contains(manifest.Path, "stack.yaml.lock") {
		lockFile := readStackLockConfigFile(manifest.Content)
		snapshop := lockFile.Snapshots[0].(map[string]interface{})["completed"]
		url := snapshop.(map[string]interface{})["url"].(string)

		for _, pkg := range lockFile.Packages {
			name, version, pkgHash, size, rev := parseExtraDep(pkg.Original.Hackage)
			if name == "" || version == "" {
				continue
			}

			metadata := HackageMetadata{
				Name:        name,
				Version:     version,
				PkgHash:     pkgHash,
				Size:        size,
				Revision:    rev,
				SnapshotURL: url,
			}

			component := types.NewComponent(name, version, Type, manifest.Path, "", metadata)
			cpes := cpe.NewCPE23(component.Name, component.Name, component.Version, Type)
			if len(cpes) > 0 {
				component.CPEs = append(component.CPEs, cpes...)
			}
			stream.AddComponent(component)
		}
	} else if strings.Contains(manifest.Path, "cabal.project.freeze") {
		packages := readManifestFile(manifest.Content)
		for _, pkg := range packages {
			name, version, pkgHash, size, rev := parseExtraDep(pkg)
			if name == "" || version == "" {
				continue
			}
			metadata := HackageMetadata{
				Name:     name,
				Version:  version,
				PkgHash:  pkgHash,
				Size:     size,
				Revision: rev,
			}

			component := types.NewComponent(name, version, Type, manifest.Path, "", metadata)
			cpes := cpe.NewCPE23(component.Name, component.Name, component.Version, Type)
			if len(cpes) > 0 {
				component.CPEs = append(component.CPEs, cpes...)
			}
			stream.AddComponent(component)
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
