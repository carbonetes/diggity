package swift

import (
	"path/filepath"
	"slices"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "swift"

var Manifests = []string{"Package.resolved", ".package.resolved"}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)
	if !ok {
		log.Error("Swift Package Manager Handler received unknown type")
	}

	metadata := readManifestFile(manifest.Content)

	switch metadata.Version {
	case 1:
		for _, pin := range metadata.Object.Pins {
			name, version := pin.Name, pin.State.Version
			component := types.NewComponent(name, version, Type, manifest.Path, "", pin)
			cpes := cpe.NewCPE23(component.Name, component.Name, component.Version, Type)
			if len(cpes) > 0 {
				component.CPEs = append(component.CPEs, cpes...)
			}
			stream.AddComponent(component)
		}
	case 2:
		for _, pin := range metadata.Pins {
			name, version := pin.Identity, pin.State.Version
			component := types.NewComponent(name, version, Type, manifest.Path, "", pin)
			cpes := cpe.NewCPE23(component.Name, component.Name, component.Version, Type)
			if len(cpes) > 0 {
				component.CPEs = append(component.CPEs, cpes...)
			}
			stream.AddComponent(component)
		}
	}

	return data
}
