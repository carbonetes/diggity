package stream

import (
	"time"

	"github.com/carbonetes/diggity/pkg/types"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

func GetComponents() []types.Component {
	data, _ := store.Get(ComponentsStoreKey)

	components, ok := data.([]types.Component)

	if !ok {
		return nil
	}

	return components
}

func GetImageInfo() types.ImageInfo {
	data, exist := store.Get(ImageInstanceStoreKey)

	if !exist {
		log.Error("ImageInstanceStore data not found")
	}

	image, ok := data.(v1.Image)

	if !ok {
		log.Error("getImageManifest received unknown data type")
	}

	digest, err := image.Digest()
	if err != nil {
		log.Error("ImageManifest digest not found :", err.Error())
	}

	mediaType, err := image.MediaType()
	if err != nil {
		log.Error("ImageManifest media type not found :", err.Error())
	}

	size, err := image.Size()
	if err != nil {
		log.Error("ImageManifest size not found :", err.Error())
	}

	manifest, err := image.Manifest()
	if err != nil {
		log.Error("ImageManifest manifest not found :", err.Error())
	}

	config, err := image.ConfigFile()
	if err != nil {
		log.Error("ImageManifest config file not found :", err.Error())
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
		log.Error("ImageManifest layers not found :", err.Error())
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
		log.Error("Distro not found")
	}

	distro, ok := data.(types.Distro)

	if !ok {
		log.Error("Invalid data type found in distro store")
	}

	return distro
}

func GetParameters() types.Parameters {
	data, exist := store.Get(ParametersStoreKey)

	if !exist {
		log.Error("Parameters not found")
	}

	parameters, ok := data.(types.Parameters)

	if !ok {
		log.Error("Invalid data type found in parameters store")
	}

	return parameters
}

func GetImageInstance() v1.Image {
	data, exist := store.Get(ImageInstanceStoreKey)

	if !exist {
		log.Error("Image Instance Store not found")
	}

	image, ok := data.(v1.Image)

	if !ok {
		log.Error("Invalid data type found in image instance store")
	}

	return image
}

func GetScanStart() time.Time {
	data, exist := store.Get(ScanStartStoreKey)

	if !exist {
		log.Error("ScanStart not found")
	}

	scanStart, ok := data.(time.Time)

	if !ok {
		log.Error("Invalid data type found in scan start store")
	}

	return scanStart
}
