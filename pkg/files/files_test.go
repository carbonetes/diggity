package files

import (
	"os"
	"testing"
)

func TestExists(t *testing.T) {
	var dir = "jenkins:2.60.3"
	if err := Exists(dir); err == true {
		t.Error("Test Failed: File exists.")
	}
}

func TestGetFilesFromDir(t *testing.T) {
	dirpath := os.TempDir()
	if _, err := os.Stat(dirpath); os.IsNotExist(err) {
		err := os.MkdirAll(dirpath, 0666)
		if err != nil{
			 t.Fatal(err.Error())
		}
	}else{
		err := os.Chmod(dirpath,0666)
		if err != nil{
			 t.Fatal(err.Error())
		}
	}
	_ , err := GetFilesFromDir(dirpath)
	
	if err != nil{
		t.Fatal(err.Error())
	}
}