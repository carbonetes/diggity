package save

import (
	"os"
	"fmt"
	"log"
	"strings"
)

// ResultToFile saves output to a file
func ResultToFile(result string, outputType *string, filename *string) {
	fileNameWithExtension := addFileExtension(*filename, *outputType)
	file, err := os.Create(fileNameWithExtension)
	if err != nil{
		log.Fatal()
	}
	defer file.Close()

	err = os.WriteFile(file.Name(), []byte(result), 0644)

	if err != nil{
		log.Fatal()
	}
	fmt.Printf("\n File Saved As : [ %v ]\n",fileNameWithExtension)
}

func addFileExtension(filename string, outputType string) string {
	removeExistingFileExtension(&filename)
	switch outputType{
	case "json" , "cdx-json", "spdx-json" :
		return filename + ".json"
	case "cdx-xml" , "spdx-xml" :
		return filename + ".xml"
	case "spdx-tag" :
		return filename + ".spdx"
	default :
		return filename + ".txt"
	}
}

// check if the filename has an existing extension
func removeExistingFileExtension(filename *string){
	currentFilename := *filename
	lastDotIndex := strings.LastIndex(currentFilename,".")

	if lastDotIndex != -1 {
		*filename = currentFilename[:lastDotIndex]
	}
}