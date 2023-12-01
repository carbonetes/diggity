package swift

import (
	"path/filepath"
	"slices"

	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "swift"


var (
	Manifests = []string{"Package.resolved", ".package.resolved"}
	log       = logger.GetLogger()
)

func CheckRelatedFile(file string) (string, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true
	}
	return "", false
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
			stream.AddComponent(component)
		}
	case 2:
		for _, pin := range metadata.Pins {
			name, version := pin.Identity, pin.State.Version
			component := types.NewComponent(name, version, Type, manifest.Path, "", pin)
			stream.AddComponent(component)
		}
	}

	return data
}
