package dpkg

import (
	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

var (
	log  = logger.GetLogger()
	Type = "dpkg"
)

func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)

	if !ok {
		log.Fatal("Dpkg received unknown file type")
	}

	attributes, err := readManifest(manifest)
	if err != nil {
		log.Error(err)
	}

	for _, attribute := range attributes {
		metadata := parseMetadata(attribute)
		if metadata == nil {
			continue
		}
		component := newComponent(*metadata)
		if len(component.Name) == 0 || len(component.Version) == 0 {
			continue
		}
		stream.AddComponent(component)
	}

	return data
}
