package stream

import (
	"time"

	"github.com/carbonetes/diggity/pkg/types"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

func GetComponents() []types.Component {
	data, exist := store.Get(ComponentsStoreKey)

	if !exist {
		log.Fatal("Components not found")
	}

	components, ok := data.([]types.Component)

	if !ok {
		log.Fatal("Invalid data type found in components store")
	}

	return components
}

func GetImageInfo() types.ImageInfo {
	data, exist := store.Get(ImageInstanceStoreKey)

	if !exist {
		return types.ImageInfo{}
	}

	image, ok := data.(v1.Image)

	if !ok {
		log.Fatal("getImageManifest received unknown data type")
	}

	digest, err := image.Digest()
	if err != nil {
		log.Fatal("ImageManifest digest not found :", err.Error())
	}

	mediaType, err := image.MediaType()
	if err != nil {
		log.Fatal("ImageManifest media type not found :", err.Error())
	}

	size, err := image.Size()
	if err != nil {
		log.Fatal("ImageManifest size not found :", err.Error())
	}

	manifest, err := image.Manifest()
	if err != nil {
		log.Fatal("ImageManifest manifest not found :", err.Error())
	}

	config, err := image.ConfigFile()
	if err != nil {
		log.Fatal("ImageManifest config file not found :", err.Error())
	}

	layers := getLayers(image)

	return types.ImageInfo{
		Digest:     digest,
		MediaType:  mediaType,
		Size:       size,
		Manifest:   *manifest,
		ConfigFile: *config,
		Layers:     layers,
	}
}

func getLayers(image v1.Image) []types.Layer {
	layers, err := image.Layers()
	if err != nil {
		log.Fatal("ImageManifest layers not found :", err.Error())
	}

	var ls []types.Layer
	for _, layer := range layers {
		digest, _ := layer.Digest()
		diffId, _ := layer.DiffID()
		size, _ := layer.Size()
		mediatype, _ := layer.MediaType()

		ls = append(ls, types.Layer{
			Digest:    digest,
			DiffID:    diffId,
			Size:      size,
			MediaType: mediatype,
		})
	}
	return ls
}

func GetDistro() types.Distro {
	data, exist := store.Get(DistroStoreKey)

	if !exist {
		log.Fatal("Distro not found")
	}

	distro, ok := data.(types.Distro)

	if !ok {
		log.Fatal("Invalid data type found in distro store")
	}

	return distro
}

func GetParameters() types.Parameters {
	data, exist := store.Get(ParametersStoreKey)

	if !exist {
		log.Fatal("Parameters not found")
	}

	parameters, ok := data.(types.Parameters)

	if !ok {
		log.Fatal("Invalid data type found in parameters store")
	}

	return parameters
}

func GetImageInstance() v1.Image {
	data, exist := store.Get(ImageInstanceStoreKey)

	if !exist {
		log.Fatal("Image Instance Store not found")
	}

	image, ok := data.(v1.Image)

	if !ok {
		log.Fatal("Invalid data type found in image instance store")
	}

	return image
}

func GetScanStart() time.Time {
	data, exist := store.Get(ScanStartStoreKey)

	if !exist {
		log.Fatal("ScanStart not found")
	}

	scanStart, ok := data.(time.Time)

	if !ok {
		log.Fatal("Invalid data type found in scan start store")
	}

	return scanStart
}

func GetParameterScanType() string {
	data, exist := store.Get(ParameterScanTypeStoreKey)

	if !exist {
		log.Fatal("ParameterScanType not found")
	}

	scanType, ok := data.(string)

	if !ok {
		log.Fatal("Invalid data type found in parameter scan type store")
	}

	return scanType
}

func GetParameterInput() string {
	data, exist := store.Get(ParameterInputStoreKey)

	if !exist {
		log.Fatal("ParameterInput not found")
	}

	input, ok := data.(string)

	if !ok {
		log.Fatal("Invalid data type found in parameter input store")
	}

	return input
}

func GetParameterOutputFormat() string {
	data, exist := store.Get(ParameterOutputFormatStoreKey)

	if !exist {
		log.Fatal("ParameterOutputFormat not found")
	}

	outputFormat, ok := data.(string)

	if !ok {
		log.Fatal("Invalid data type found in parameter output format store")
	}

	return outputFormat
}

func GetParameterQuiet() bool {
	data, exist := store.Get(ParameterQuietStoreKey)

	if !exist {
		log.Fatal("ParameterQuiet not found")
	}

	quiet, ok := data.(bool)

	if !ok {
		log.Fatal("Invalid data type found in parameter quiet store")
	}

	return quiet
}

func GetParameterMaxFileSize() int64 {
	data, exist := store.Get(ParameterMaxFileSizeStoreKey)

	if !exist {
		log.Fatal("ParameterMaxFileSize not found")
	}

	maxFileSize, ok := data.(int64)

	if !ok {
		log.Fatal("Invalid data type found in parameter max file size store")
	}

	return maxFileSize
}

func GetParameterScanners() []string {
	data, exist := store.Get(ParameterScannersStoreKey)

	if !exist {
		log.Fatal("ParameterScanners not found")
	}

	scanners, ok := data.([]string)

	if !ok {
		log.Fatal("Invalid data type found in parameter scanners store")
	}

	return scanners
}

func GetParameterAllowFileListing() bool {
	data, exist := store.Get(ParameterAllowFileListingStoreKey)

	if !exist {
		log.Fatal("ParameterAllowFileListing not found")
	}

	allowFileListing, ok := data.(bool)

	if !ok {
		log.Fatal("Invalid data type found in parameter allow file listing store")
	}

	return allowFileListing
}

func GetParameterRegistry() string {
	data, exist := store.Get(ParameterRegistryStoreKey)

	if !exist {
		log.Fatal("ParameterRegistry not found")
	}

	registry, ok := data.(string)

	if !ok {
		log.Fatal("Invalid data type found in parameter registry store")
	}

	return registry
}

func GetFiles() []string {
	data, exist := store.Get(FileListStoreKey)

	if !exist {
		log.Fatal("FileList not found")
	}

	files, ok := data.([]string)

	if !ok {
		log.Fatal("Invalid data type found in file list store")
	}

	return files
}

func CheckIfComponentExists(component types.Component) bool {
	components := GetComponents()

	for _, c := range components {
		if c.Name == component.Name && c.Version == component.Version && c.Type == component.Type {
			return true
		}
	}
	return false
}