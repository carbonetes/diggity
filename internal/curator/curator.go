package curator

import (
	"fmt"
	"log"
)

func IndexImageFilesystem(data interface{}) interface{} {
	imageName, ok := data.(string)
	if !ok {
		log.Fatal("IndexImageFilesystem received unknown type")
	}
	fmt.Println(imageName)
	image, err := GetImage(imageName)
	if err != nil {
		log.Fatal(err)
	}

	err = ReadFiles(image)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
