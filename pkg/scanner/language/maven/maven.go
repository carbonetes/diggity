package maven

import (
	"path/filepath"
	"slices"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/helper"
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
	payload, ok := data.(types.Payload)
	if !ok {
		log.Error("Java Archive received unknown file type")
		return nil
	}

	scan(payload)

	return data
}

func scan(payload types.Payload) {
	manifest := payload.Body.(types.ManifestFile)
	metadata := readManifestFile(manifest.Content)
	if metadata == nil {
		return
	}

	if len(metadata.Dependencies) > 0 {
		for _, dependency := range metadata.Dependencies {
			if dependency.GroupID == "" || dependency.ArtifactID == "" || dependency.Version == "" {
				continue
			}

			c := component.New(dependency.ArtifactID, dependency.Version, Type)

			cpes := cpe.NewCPE23(c.Name, c.Name, c.Version, Type)
			if len(cpes) > 0 {
				for _, cpe := range cpes {
					component.AddCPE(c, cpe)
				}
			}

			// Correction for PackageURL
			c.PackageURL = "pkg:maven/" + dependency.GroupID + "/" + dependency.ArtifactID + "@" + dependency.Version

			component.AddOrigin(c, manifest.Path)
			component.AddType(c, Type)

			rawMetadata, err := helper.ToJSON(dependency)
			if err != nil {
				log.Errorf("Failed to convert metadata to JSON: %v", err)
			}

			if len(rawMetadata) > 0 {
				component.AddRawMetadata(c, rawMetadata)
			}

			cdx.AddComponent(c, payload.Address)
		}
	}

	if metadata.ArtifactID == "" || metadata.Version == "" {
		return
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

	rawMetadata, err := helper.ToJSON(metadata)
	if err != nil {
		log.Errorf("Failed to convert metadata to JSON: %v", err)
	}

	if len(rawMetadata) > 0 {
		component.AddRawMetadata(c, rawMetadata)
	}

	cdx.AddComponent(c, payload.Address)

}
