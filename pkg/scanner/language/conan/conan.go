package conan

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

const Type string = "conan"

var Manifests = []string{"conanfile.txt", "conan.lock"}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	payload, ok := data.(types.Payload)
	if !ok {
		log.Error("Conan Handler received unknown type")
		return nil
	}

	scan(payload)

	return data
}

func scan(payload types.Payload) {
	manifest := payload.Body.(types.ManifestFile)
	if strings.Contains(manifest.Path, "conanfile.txt") {
		packages := readManifestFile(manifest.Content)
		if len(packages) == 0 {
			return
		}

		for _, pkg := range packages {
			if pkg.Name == "" {
				continue
			}

			c := component.New(pkg.Name, pkg.Version, Type)
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

			if len(payload.Layer) > 0 {
				component.AddLayer(c, payload.Layer)
			}

			cdx.AddComponent(c, payload.Address)
		}
	} else if strings.Contains(manifest.Path, "conan.lock") {
		metadata := readLockFile(manifest.Content)
		if len(metadata.GraphLock.Nodes) == 0 {
			return
		}
		for _, node := range metadata.GraphLock.Nodes {
			if node.Ref == "" {
				continue
			}
			var val []string
			if strings.Contains(node.Ref, "@") {
				val = strings.Split(node.Ref, "@")
			} else if strings.Contains(node.Ref, "#") {
				val = strings.Split(node.Ref, "#")
			} else {
				val = []string{node.Ref}
			}

			ref := strings.Split(val[0], "/")
			if len(ref) != 2 {
				continue
			}
			name, version := ref[0], ref[1]

			if strings.Contains(version, "[") && strings.Contains(version, "]") {
				version = strings.Replace(version, "[", "", -1)
				version = strings.Replace(version, "]", "", -1)
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

			rawMetadata, err := helper.ToJSON(node)
			if err != nil {
				log.Errorf("Error converting metadata to JSON: %s", err)
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
}
