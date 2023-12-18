package gradle

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/stream"
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
	}

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

		component := types.NewComponent(metadata.Name, metadata.Version, Type, manifest.Path, "", metadata)

		stream.AddComponent(component)
	}
	return data
}
