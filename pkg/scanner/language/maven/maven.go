package maven

import (
	"path/filepath"
	"slices"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "java"

var Manifests = []string{"pom.xml", "pom.properties"}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)
	if !ok {
		log.Error("Java Archive received unknown file type")
		return nil
	}

	metadata := readManifestFile(manifest.Content)
	if metadata == nil {
		return nil
	}

	if metadata.ArtifactID == "" || metadata.Version == "" {
		return nil
	}

	component := types.NewComponent(metadata.ArtifactID, metadata.Version, Type, manifest.Path, metadata.Description, metadata)
	component.PURL = "pkg:maven/" + metadata.GroupID + "/" + metadata.ArtifactID + "@" + metadata.Version
	cpes := cpe.NewCPE23(component.Name, component.Name, component.Version, Type)
	if len(cpes) > 0 {
		component.CPEs = append(component.CPEs, cpes...)
	}

	stream.AddComponent(component)

	return data
}