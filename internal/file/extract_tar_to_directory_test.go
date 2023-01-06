package file

import (
	"path/filepath"
	"testing"
)

func TestUntar(t *testing.T) {
	var dst = filepath.Join("C:", "Users", "Username", "AppData", "Local", "Temp", "1748738324", "diggity-tmp-f72ce52d-1cc9-47d4-b672-d94e8a7ee78e", "0a9bab4e2f50283e7061494edf84f9e8695fd75813fe55dc2d11498d1996f397")
	var source = filepath.Join("C:", "Users", "Username", "AppData", "Local", "Temp", "1748738324", "diggity-tmp-f72ce52d-1cc9-47d4-b672-d94e8a7ee78e", "0a9bab4e2f50283e7061494edf84f9e8695fd75813fe55dc2d11498d1996f397", "layer.tar")
	var recursive = true
	err := UnTar(dst, source, recursive)
	if err != nil && !recursive {
		t.Error("Testing Failed: An error occured when processing.")
	}
}
