package maven

import (
	"path/filepath"
	"slices"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
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

	c := component.New(metadata.ArtifactID, metadata.Version, Type)

	cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
	if len(cpes) > 0 {
		for _, cpe := range cpes {
			component.AddCPE(c, cpe)
		}
	}

	// Correction for PackageURL
	c.PackageURL = "pkg:maven/" + metadata.GroupID + "/" + metadata.ArtifactID + "@" + metadata.Version

	component.AddOrigin(c, manifest.Path)
	component.AddType(c, Type)
	component.AddDescription(c, metadata.Description)

	cdx.AddComponent(c)

	return data
}
