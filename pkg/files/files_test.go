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
	filePath := "../../testFile.txt"

	file, err := os.Create(filePath)
	if err != nil{
		t.Fatal(err.Error())
	}
	file.Close()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err := os.MkdirAll(filePath, 0666)
		if err != nil{
			 t.Fatal(err.Error())
		}
	}else{
		err := os.Chmod(filePath,0666)
		if err != nil{
			 t.Fatal(err.Error())
		}
	}

	_ , err = GetFilesFromDir(filePath)
	if err != nil{
		t.Fatal(err.Error())
	}

	err = os.Remove(filePath)
	if err != nil{
		t.Fatal(err.Error())
	}
}