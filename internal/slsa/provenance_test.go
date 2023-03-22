package slsa

import (
	"path/filepath"
	"testing"
)

func TestGetProvenanceMetadata(t *testing.T) {
	path := filepath.Join("..", "..", "docs", "references", "provenance", "build.provenance")
	output, err := GetProvenanceMetadata(path)
	if err != nil {
		t.Error("Test Failed: Error occurred parsing provenance file.")
	}
	if output == nil {
		t.Error("Test Failed: No provenance metadata parsed.")
	}
}
