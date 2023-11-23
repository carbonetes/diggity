package javaarchive

import (
	"github.com/carbonetes/diggity/pkg/types"
	"github.com/google/uuid"
)

func newComponent(pom types.JavaManifest) types.Component {
	return types.Component{
		ID:       uuid.New().String(),
		Name:     pom.ArtifactID,
		Version:  pom.Version,
		Type:     Type,
		Metadata: pom,
	}
}
