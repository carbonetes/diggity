package stream

import (
	"log"

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

func GetImageManifest() types.ImageManifest {
	data, exist := store.Get(ImageInstanceStoreKey)

	if !exist {
		log.Fatal("ImageInstanceStore data not found")
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

	return types.ImageManifest{
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

func GetSecretParameters() types.SecretParameters {
	data, exist := store.Get(SecretParametersStoreKey)

	if !exist {
		log.Fatal("SecretParameters not found")
	}

	parameters, ok := data.(types.SecretParameters)

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
