package slsa

import (
	"path/filepath"
	"testing"
)

type (
	ValidateProvenanceResult struct {
		input    map[string]interface{}
		expected bool
	}
)

var (
	testProvenance = map[string]interface{}{
		"_type":         "https://in-toto.io/Statement/v0.1",
		"predicate":     "predicateValue",
		"predicateType": "https://slsa.dev/provenance/v0.1",
		"subject":       "subjectValue",
	}
	invalidProvenance = map[string]interface{}{
		"key": "val",
	}
)

func TestParseProvenanceMetadata(t *testing.T) {
	path := filepath.Join("..", "..", "docs", "references", "provenance", "build.provenance")
	output, err := parseProvenanceMetadata(path)
	if err != nil {
		t.Error("Test Failed: Error occurred parsing provenance file.")
	}
	if output == nil {
		t.Error("Test Failed: No provenance metadata parsed.")
	}
}

func TestValidateProvenance(t *testing.T) {
	tests := []ValidateProvenanceResult{
		{testProvenance, true},
		{invalidProvenance, false},
		{nil, false},
	}

	for _, test := range tests {
		if output := validateProvenance(test.input); output != test.expected {
			t.Error("Test Failed: Boolean results must be the same.")
		}
	}
}
