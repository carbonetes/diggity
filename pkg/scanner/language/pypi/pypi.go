package pypi

import (
	"path/filepath"
	"slices"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "pypi"

var (
	Manifests  = []string{"METADATA", "requirements.txt", "poetry.lock"}
	Extensions = []string{".egg-info"}
)

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) || slices.Contains(Extensions, filepath.Ext(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)
	if !ok {
		log.Error("Python Handler received unknown type")
	}

	if filepath.Ext(manifest.Path) == ".egg-info" || filepath.Base(manifest.Path) == "METADATA" {
		metadata := readManifestFile(manifest.Content)
		name, version := metadata["Name"].(string), metadata["Version"].(string)
		component := types.NewComponent(name, version, Type, manifest.Path, "", metadata)
		if val, ok := metadata["Summary"].(string); ok {
			component.Description = val
		}
		if val, ok := metadata["License"].(string); ok {
			component.Licenses = append(component.Licenses, val)
		}
		cpes := cpe.NewCPE23(component.Name, component.Name, component.Version, Type)
		if len(cpes) > 0 {
			component.CPEs = append(component.CPEs, cpes...)
		}
		stream.AddComponent(component)
	} else if filepath.Base(manifest.Path) == "requirements.txt" {
		attributes := readRequirementsFile(manifest.Content)
		for _, attribute := range attributes {
			name, version := attribute[0], attribute[1]
			metadata := map[string]string{"name": name, "version": version}
			component := types.NewComponent(name, version, Type, manifest.Path, "", metadata)
			cpes := cpe.NewCPE23(component.Name, component.Name, component.Version, Type)
			if len(cpes) > 0 {
				component.CPEs = append(component.CPEs, cpes...)
			}
			stream.AddComponent(component)
		}
	} else if filepath.Base(manifest.Path) == "poetry.lock" {
		metadata := readPoetryLockFile(manifest.Content)
		for _, packageInfo := range metadata.Packages {
			name, version := packageInfo.Name, packageInfo.Version
			component := types.NewComponent(name, version, Type, manifest.Path, "", packageInfo)
			cpes := cpe.NewCPE23(component.Name, component.Name, component.Version, Type)
			if len(cpes) > 0 {
				component.CPEs = append(component.CPEs, cpes...)
			}
			stream.AddComponent(component)
		}
	}

	return data
}
