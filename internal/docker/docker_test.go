package docker

import (
	"testing"
)

func TestConnectionTest(t *testing.T) {
	if err := testConnection(); err != nil {
		t.Error("Test Failed: Error Occurred upon testing connection.")
	}
}
