package file

import (
	"testing"
)

func TestExists(t *testing.T) {
	var dir = "jenkins:2.60.3"
	if err := Exists(dir); err == true {
		t.Error("Test Failed: File exists.")
	}
}
