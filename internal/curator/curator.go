package curator

import (
	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/pkg/stream"
)

var log = logger.GetLogger()

func IndexImageFilesystem(data interface{}) interface{} {
	imageName, ok := data.(string)
	if !ok {
		log.Error("IndexImageFilesystem received unknown type")
	}
	image, err := GetImage(imageName)
	if err != nil {
		log.Error(err)
	}
	stream.SetImageInstance(image)
	err = ReadFiles(image)
	if err != nil {
		log.Error(err)
	}
	return data
}

func IndexTarballFilesystem(data interface{}) interface{} {
	tarballPath, ok := data.(string)
	if !ok {
		log.Error("IndexTarballFilesystem received unknown type")
	}
	image, err := ReadTarball(tarballPath)
	if err != nil {
		log.Error(err)
	}
	stream.SetImageInstance(image)
	err = ReadFiles(image)
	if err != nil {
		log.Error(err)
	}
	return data
}
