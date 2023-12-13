package conan

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/stream"
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
	manifest, ok := data.(types.ManifestFile)
	if !ok {
		log.Error("Conan Handler received unknown type")
	}

	if strings.Contains(manifest.Path, "conanfile.txt") {
		packages := readManifestFile(manifest.Content)
		if len(packages) == 0 {
			return nil
		}

		for _, pkg := range packages {
			if pkg.Name == "" {
				continue
			}

			component := types.NewComponent(pkg.Name, pkg.Version, Type, manifest.Path, "", pkg)
			stream.AddComponent(component)
		}
	} else if strings.Contains(manifest.Path, "conan.lock") {
		metadata := readLockFile(manifest.Content)
		if len(metadata.GraphLock.Nodes) == 0 {
			return nil
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

			component := types.NewComponent(name, version, Type, manifest.Path, "", node)
			stream.AddComponent(component)
		}
	}

	return data
}
