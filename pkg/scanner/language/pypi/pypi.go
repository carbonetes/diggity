package pypi

import (
	"path/filepath"
	"slices"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "pypi"

var (
	Manifests  = []string{"METADATA", "requirements.txt", "poetry.lock", "PKG-INFO"}
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
		return nil
	}

	scan(manifest)

	return data
}

func scan(manifest types.ManifestFile) {
	if filepath.Ext(manifest.Path) == ".egg-info" || filepath.Base(manifest.Path) == "METADATA" || filepath.Base(manifest.Path) == "PKG-INFO" {
		metadata := readManifestFile(manifest.Content)
		name, version := metadata["Name"].(string), metadata["Version"].(string)

		c := component.New(name, version, Type)

		cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
		if len(cpes) > 0 {
			for _, cpe := range cpes {
				component.AddCPE(c, cpe)
			}
		}

		component.AddOrigin(c, manifest.Path)
		component.AddType(c, Type)

		if val, ok := metadata["Summary"].(string); ok {
			component.AddDescription(c, val)
		}

		if val, ok := metadata["License"].(string); ok {
			component.AddLicense(c, val)
		}

		rawMetadata, err := helper.ToJSON(metadata)
		if err != nil {
			log.Errorf("Error converting metadata to JSON: %s", err)
		}

		if len(rawMetadata) > 0 {
			component.AddRawMetadata(c, rawMetadata)
		}

		cdx.AddComponent(c)

	} else if filepath.Base(manifest.Path) == "requirements.txt" {
		attributes := readRequirementsFile(manifest.Content)
		for _, attribute := range attributes {
			name, version := attribute[0], attribute[1]

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
	} else if filepath.Base(manifest.Path) == "poetry.lock" {
		metadata := readPoetryLockFile(manifest.Content)
		for _, packageInfo := range metadata.Packages {
			name, version := packageInfo.Name, packageInfo.Version

			c := component.New(name, version, Type)

			cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
			if len(cpes) > 0 {
				for _, cpe := range cpes {
					component.AddCPE(c, cpe)
				}
			}

			component.AddOrigin(c, manifest.Path)
			component.AddType(c, Type)

			rawMetadata, err := helper.ToJSON(packageInfo)
			if err != nil {
				log.Errorf("Error converting metadata to JSON: %s", err)
			}

			if len(rawMetadata) > 0 {
				component.AddRawMetadata(c, rawMetadata)
			}

			cdx.AddComponent(c)
		}
	}
}
