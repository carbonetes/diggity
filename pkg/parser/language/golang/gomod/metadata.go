package gomod

import (
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"golang.org/x/mod/modfile"
)

// Initialize go metadata values from content
func parseMetadata(modPkg interface{}) *metadata.GoModMetadata {
	var metadata metadata.GoModMetadata

	switch goPkg := modPkg.(type) {
	case *modfile.Require:
		metadata.Path = goPkg.Mod.Path
		metadata.Version = goPkg.Mod.Version
	case *modfile.Replace:
		metadata.Path = goPkg.New.Path
		metadata.Version = goPkg.New.Version
	}
	return &metadata
}
