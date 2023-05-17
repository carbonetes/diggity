package docker

import (
	"path/filepath"
	"testing"
)

const dockerReference string = "diggity-tmp-385abb3c-df38-44dd-b30f-467ba364ee3a"

func TestGetJSONFilesFromDir(t *testing.T) {
	
	rootPath := filepath.Join("..", "..", "..", "docs", "references", "docker", dockerReference)
	expected := []string{
		filepath.Join("..", "..", "..", "docs", "references", "docker", dockerReference, "bfe296a525011f7eb76075d688c681ca4feaad5afe3b142b36e30f1a171dc99a.json"),
		filepath.Join("..", "..", "..", "docs", "references", "docker", dockerReference, "manifest.json"),
	}

	output, err := getJSONFilesFromDir(rootPath)

	if err != nil {
		t.Error("Test Failed: Error occurred while parsing docker files.")
	}

	if len(output) != len(expected) {
		t.Errorf("Test Failed: Slice length must be equal with the expected result. Expected: %v, Received: %v", len(expected), len(output))
	}

	for i, file := range output {
		if file != expected[i] {
			t.Errorf("Test Failed: Expected output of %v, received: %v ", expected[i], file)
		}
	}
}
