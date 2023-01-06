package docker

import (
	"path/filepath"
	"testing"

	"github.com/docker/docker/pkg/ioutils"
)

func TestDir(t *testing.T) {
	tests := []string{
		filepath.Join("dir", "test", "123456789"),
		filepath.Join("C:", "Users", "Username", "AppData", "Local", "Temp", "1156647128"),
		filepath.Join("Library", "Caches", "TemporaryItems", "778637836"),
	}

	for _, test := range tests {
		tempDir = test
		if output := Dir(); output != tempDir {
			t.Errorf("Test Failed: Expected output of %v, received: %v", tempDir, output)
		}
	}
}

func TestExtractedDir(t *testing.T) {
	tests := []string{
		"test/dir",
		filepath.Join("test", "path", "dir"),
		filepath.Join("C:", "Users", "Username", "AppData", "Local", "Temp", "3813774253", "diggity-tmp-865bf4ad-4547-4acb-9cd0-98b48e8b347e"),
		"",
	}

	for _, test := range tests {
		extractDir = test
		if output := ExtractedDir(); output != extractDir {
			t.Errorf("Test Failed: Expected output of %v, received: %v", tempDir, output)
		}
	}
}

func TestConnectionTest(t *testing.T) {
	if err := testConnection(); err != nil {
		t.Error("Test Failed: Error Occurred upon testing connection.")
	}
}

func TestCreateTempDir(t *testing.T) {
	CreateTempDir()
	expected, _ := ioutils.TempDir("", "")
	if filepath.Dir(tempDir) != filepath.Dir(expected) {
		t.Errorf("Test Failed: Expected output of %v, received: %v", filepath.Dir(expected), filepath.Dir(tempDir))
	}
}
