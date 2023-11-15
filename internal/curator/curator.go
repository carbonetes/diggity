package curator

import (
	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/pkg/stream"
)

var log = logger.GetLogger()

// ImageScanHandler scans the given image for files and reads them.
// It takes in a string parameter as the image name and returns the same parameter.
func ImageScanHandler(data interface{}) interface{} {
	imageName, ok := data.(string)
	if !ok {
		log.Fatal("IndexImageFilesystem received unknown type")
	}
	image, err := GetImage(imageName)
	if err != nil {
		log.Fatal(err)
	}
	stream.SetImageInstance(image)
	err = ReadFiles(image)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

// TarballScanHandler scans a tarball file and reads its contents as an image instance.
// It then sets the image instance to the stream and reads the files in the image.
func TarballScanHandler(data interface{}) interface{} {
	tarballPath, ok := data.(string)
	if !ok {
		log.Error("TarballStoreWatcher received unknown type")
	}
	image, err := ReadTarballAsImage(tarballPath)
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
