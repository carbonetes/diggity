package apk

import (
	"slices"
	"strings"

	"github.com/carbonetes/diggity/pkg/types"
)

var Manifests = []string{"lib/apk/db/installed"}

func readManifest(manifest types.ManifestFile) ([]string, error) {
	packages := strings.Split(string(manifest.Content), "\n\n")

	return packages, nil
}

func CheckRelatedFile(file string) (string, bool) {
	if slices.Contains(Manifests, file) {
		return Type, true
	}
	return "", false
}
