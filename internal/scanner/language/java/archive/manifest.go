package javaarchive

import (
	"path/filepath"
	"slices"
)

var Manifests = []string{"pom.xml", "pom.properties"}

func CheckRelatedFile(file string) (string, bool) {
	if slices.Contains(Manifests, filepath.Base(file)) {
		return Type, true
	}
	return "", false
}
