package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExists(t *testing.T) {
	var dir = "jenkins:2.60.3"
	if result := Exists(dir); result == true {
		t.Error("Test Failed: File exists.")
	}
}

func TestGetFilesFromDir(t *testing.T) {
	filePath := filepath.Join(os.TempDir(), "test_files")

	file, err := os.Create(filePath)
	if err != nil {
		t.Fatal(err.Error())
	}
	file.Close()

	_, err = os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Error("Test Failed: File does not exist!")
			t.Fail()
		} else {
			t.Error(err.Error())
			t.Fail()
		}
	}

	result, err := GetFilesFromDir(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Error("Test Failed: File does not exist!")
			t.Fail()
		} else {
			t.Error(err.Error())
			t.Fail()
		}
	}

	if result == nil || len(*result) == 0 {
		t.Error("Result is nil")
		t.Fail()
	}

	err = os.Remove(filePath)
	if err != nil {
		t.Error(err.Error())
	}
}
