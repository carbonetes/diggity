package javaarchive

import (
	"encoding/xml"
	"io"

	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type = "java-archive"

var log = logger.GetLogger()

func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)
	if !ok {
		log.Fatal("Java Archive received unknown file type")
		return nil
	}

	var pom types.JavaManifest

	err := xml.Unmarshal(manifest.Content, &pom)
	if err != nil {
		if err == io.EOF {
			return nil
		} else {
			log.Error(err)
			return nil
		}
	}

	component := newComponent(pom)
	if len(component.Name) == 0 || len(component.Version) == 0 {
		return data
	}

	stream.AddComponent(component)

	return data
}
