package dpkg

import (
	"slices"

	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/pkg/types"
)

var Manifests = []string{"var/lib/dpkg/status"}

func CheckRelatedFile(file string) (string, bool) {
	if slices.Contains(Manifests, file) {
		return Type, true
	}
	return "", false
}

func readManifest(manifest types.ManifestFile) ([]string, error) {
	contents := string(manifest.Content)
	packages := helper.SplitContentsByEmptyLine(contents)

	return packages, nil
}
