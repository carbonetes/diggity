package cocoapods

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "cocoapods"

var Manifests = []string{"Podfile.lock"}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)
	if !ok {
		log.Error("Cocoapods Handler received unknown type")
	}

	metadata := readManifestFile(manifest.Content)
	for _, pod := range metadata.Pods {
		var pods string
		switch c := pod.(type) {
		case string:
			pods = c
		case map[string]interface{}:
			val := pod.(map[string]interface{})
			for all := range val {
				pods = all
			}
		}

		attributes := strings.Split(pods, " ")
		name, version := attributes[0], strings.TrimSuffix(strings.TrimPrefix(attributes[1], "("), ")")
		base := strings.Split(name, "/")[0]
		component := types.NewComponent(name, version, Type, manifest.Path, "", nil)
		if val, ok := metadata.SpecChecksums[base]; ok {
			component.Metadata = FileLockMetadataCheckSums{Checksums: val}
		}
		cpes := cpe.NewCPE23(component.Name, component.Name, component.Version, Type)
		if len(cpes) > 0 {
			component.CPEs = append(component.CPEs, cpes...)
		}

		stream.AddComponent(component)
	}

	return data
}
