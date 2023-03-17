package secret

import "testing"

type (
	validateFilenameResult struct {
		input     string
		reference map[string]string
		expected  bool
	}
)

var (
	extensions = map[string]string{
		".env":        ".env",
		".h":          ".h",
		".so":         ".so",
		".pem":        ".pem",
		".properties": ".properties",
		".xml":        ".xml",
		".yml":        ".yml",
		".yaml":       ".yaml",
		".json":       ".json",
		".py":         ".py",
		".js":         ".js",
		".ts":         ".ts",
		".PHP":        ".PHP",
	}
)

func TestValidateFilename(t *testing.T) {
	tests := []validateFilenameResult{
		{"invalid.tar", extensions, false},
		{"invalid.gz", extensions, false},
		{"test.env", extensions, true},
		{"test.h", extensions, true},
		{"test.so", extensions, true},
		{"test", extensions, true},
		{"test.xyz", extensions, false},
		{"test.null", extensions, false},
	}

	for _, test := range tests {
		if output := validateFilename(test.input, test.reference); output != test.expected {
			t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected, output)
		}
	}
}
