package relationship

import (
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/docker"
	"github.com/carbonetes/diggity/pkg/parser/source"
)

const (
	contains    = "contains"
	layerSuffix = "/layer.tar"
)

// Find packages contained by the source
func FindSourceContains() {
	var parentLayer string

	if *bom.Arguments.Dir == "" {
		layers := docker.ImageInfo.DockerManifest[0].Layers.([]interface{})
		layer := layers[0].(string)
		parentLayer = strings.Replace(layer, layerSuffix, "", -1)
	} else {
		parentLayer = source.SourceInfo.ID
	}

	for _, pkg := range bom.Packages {
		Relationships = append(Relationships, model.Relationship{
			Parent: parentLayer,
			Child:  pkg.ID,
			Type:   contains,
		})
	}

	wg.Done()
}
