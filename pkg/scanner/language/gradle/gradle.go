package gradle

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/component"
	"github.com/carbonetes/diggity/pkg/types"
)

const Type string = "gradle"

var Manifests = []string{"buildscript-gradle.lockfile", ".build.gradle"}

func CheckRelatedFile(file string) (string, bool, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true, true
	}
	return "", false, false
}

func Scan(data interface{}) interface{} {
	manifest, ok := data.(types.ManifestFile)
	if !ok {
		log.Error("Gradle Handler received unknown type")
		return nil
	}

	scan(manifest)

	return data
}

func scan(manifest types.ManifestFile) {
	lines := strings.Split(string(manifest.Content), "\n")
	for _, line := range lines {
		if !strings.Contains(line, ":") {
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		attributes := strings.SplitN(line, ":", 3)
		if len(attributes) < 3 {
			continue
		}

		metadata := Metadata{
			Vendor:  attributes[0],
			Name:    attributes[1],
			Version: strings.ReplaceAll(attributes[2], "=classpath", ""),
		}

		c := component.New(metadata.Name, metadata.Version, Type)

		cpes := cpe.NewCPE23(metadata.Vendor, c.Name, c.Version, Type)
		if len(cpes) > 0 {
			for _, cpe := range cpes {
				component.AddCPE(c, cpe)
			}
		}

		component.AddOrigin(c, manifest.Path)
		component.AddType(c, Type)

		rawMetadata, err := helper.ToJSON(metadata)
		if err != nil {
			log.Errorf("Error converting metadata to JSON: %s", err)
		}

		if len(rawMetadata) > 0 {
			component.AddRawMetadata(c, rawMetadata)
		}

		cdx.AddComponent(c)
	}
}
