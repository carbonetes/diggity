package pub

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "pub"

var Manifests = []string{"pubspec.yaml", "pubspec.lock"}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)
	if !ok {
		log.Error("Pub Handler received unknown type")
		return nil
	}

	scan(manifest)

	return data
}

func scan(manifest types.ManifestFile) {
	if strings.Contains(manifest.Path, "pubspec.yaml") {
		metadata := readManifestFile(manifest.Content)
		var name, version, license string

		name, ok := metadata["name"].(string)
		if !ok {
			return
		}
		version, ok = metadata["version"].(string)
		if !ok {
			version = "0.0.0"
		}

		if val, ok := metadata["license"].(string); ok {
			license = val
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

		rawMetadata, err := helper.ToJSON(metadata)
		if err != nil {
			log.Errorf("Error converting metadata to JSON: %s", err)
		}

		if len(rawMetadata) > 0 {
			component.AddRawMetadata(c, rawMetadata)
		}

		if license != "" {
			component.AddLicense(c, license)
		}

		cdx.AddComponent(c)

	} else if strings.Contains(manifest.Path, "pubspec.lock") {
		metadata := readLockFile(manifest.Content)
		if len(metadata.Packages) == 0 {
			return
		}
		for _, pkg := range metadata.Packages {
			if pkg.Description.Name == "" || pkg.Version == "" {
				continue
			}

			c := component.New(pkg.Description.Name, pkg.Version, Type)

			cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
			if len(cpes) > 0 {
				for _, cpe := range cpes {
					component.AddCPE(c, cpe)
				}
			}

			component.AddOrigin(c, manifest.Path)
			component.AddType(c, Type)

			rawMetadata, err := helper.ToJSON(pkg)
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
