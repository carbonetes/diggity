package apk

import (
	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

const Distro = "alpine"

var (
	log = logger.GetLogger()
	Type = "apk"
)

func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)

	if !ok {
		log.Error("Apk Handler received unknown type")
		return nil
	}

	attributes, err := readManifest(manifest)
	if err != nil {
		log.Error(err)
	}

	for _, attribute := range attributes {
		component := newComponent(attribute)
		if len(component.ID) == 0 {
			continue
		}
		stream.AddComponent(component)
	}
	return data
}
