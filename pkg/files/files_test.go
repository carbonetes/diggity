package files

import (
	"os"
	"strings"
	"testing"
)

func TestExists(t *testing.T) {
	var dir = "jenkins:2.60.3"
	if err := Exists(dir); err == true {
		t.Error("Test Failed: File exists.")
	}
}

func TestGetFilesFromDir(t *testing.T) {
	currentDirectory := os.TempDir()
	newDirectory := strings.Replace(currentDirectory, "Temp", ".", -1)

	if _, err := os.Stat(newDirectory); os.IsNotExist(err) {
		err := os.MkdirAll(newDirectory, 0666)
		if err != nil{
			 t.Fatal(err.Error())
		}
	}else{
		err := os.Chmod(newDirectory,0666)
		if err != nil{
			 t.Fatal(err.Error())
		}
	}
	_ , err := GetFilesFromDir(newDirectory)
	
	if err != nil{
		t.Fatal(err.Error())
	}
}