package bom

import (
	"testing"

	"github.com/carbonetes/diggity/pkg/docker"
	"github.com/carbonetes/diggity/pkg/provider"
	"github.com/stretchr/testify/assert"
)

func TestInitParsers(t *testing.T) {

	// Test case 1: Image argument provided
	arg1 := provider.NewArguments()
	arg1.Image = stringPtr("alpine")

	InitParsers(*arg1)
	if !assert.DirExists(t, *Target) {
		t.Errorf("Target was not set correctly '%s'", *Target)
	}

	// Test case 2: Dir argument provided
	arg2 := provider.NewArguments()
	arg2.Dir = stringPtr(".")

	InitParsers(*arg2)
	if !assert.DirExists(t, *Target) {
		t.Errorf("Target was not set correctly '%s'", *Target)
	}

	// Test case 3: Tar argument provided
	tarFile := docker.SaveImageToTar(stringPtr("alpine"))
	arg3 := provider.NewArguments()
	arg3.Tar = stringPtr(tarFile.Name())

	InitParsers(*arg3)
	if !assert.DirExists(t, *Target) {
		t.Errorf("Target was not set correctly '%s'", *Target)
	}
}

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}
