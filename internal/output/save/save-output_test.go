package save

import (
	"os"
	"testing"
)

type validateFilenameExtension struct{
	filename 	string
	outputType  string
	expected 	string
}

func TestAddFileExtension(t *testing.T) {
	tests:= []validateFilenameExtension{
		{"result.file.json","json", "result.file.json"},
		{"result.json","json", "result.json"},
		{"result","json", "result.json"},
		{"result","cyclonedx-json", "result.json"},
		{"result","spdx-json", "result.json"},
		{"result","cyclonedx-xml", "result.xml"},
		{"result","spdx-xml", "result.xml"},
		{"result","spdx-tag-value", "result.spdx"},
		{"result","table", "result.txt"},
	}
 
	for _, test := range tests {
		if result := addFileExtension(test.filename,test.outputType); result != test.expected{
			t.Errorf(" Test Failed: Expected output of %v , Received: %v ", test.expected, result);
		}
	}
}

func TestResultToFile(t *testing.T) {
	filename := "result"
	outputType := "txt"
	fileContent := `BOM Diggity's primary purpose is to ensure the security and integrity of software programs. It incorporates secret analysis allowing the user to secure crucial information before deploying any parts of the application to the public.`
	filenameWithExtension := filename + "." + outputType

	ResultToFile(fileContent, &outputType, &filename)
	stat , err := os.Stat(filenameWithExtension)
	if os.IsNotExist(err) {
		t.Fatalf("file %s was not created", filenameWithExtension)
	}

	data , err := os.ReadFile(stat.Name())
	if err != nil {
		t.Fatal("Error reading File")
	}

	if string(data) != fileContent {
		t.Fatal("The File has incorrect content")
	}
 
	err = os.Remove(stat.Name())
	if err != nil{
		t.Fatal(err.Error())
	}
}
