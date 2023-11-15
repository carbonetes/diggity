package stream

import (
	"github.com/carbonetes/diggity/pkg/types"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

func SetDefaultValues() {
	store.Set(ComponentsStoreKey, []types.Component{})
	store.Set(SecretsStoreKey, []types.Secret{})
	store.Set(SBOMStoreKey, types.NewSBOM())
	store.Set(DistroStoreKey, types.Distro{})
	store.Set(ParametersStoreKey, types.Parameters{})
}

func AddComponent(component types.Component) {
	data, exist := store.Get(ComponentsStoreKey)

	components, ok := data.([]types.Component)

	if !ok {
		log.Error("Received invalid component slice from store")
	}

	if !exist {
		store.Set(ComponentsStoreKey, []types.Component{component})
	}

	components = append(components, component)
	store.Set(ComponentsStoreKey, components)
}

func AddSecret(secret types.Secret) {
	data, exist := store.Get(SecretsStoreKey)

	secrets, ok := data.([]types.Secret)

	if !ok {
		log.Error("Received invalid secret slice from store")
	}

	if !exist {
		store.Set(SecretsStoreKey, []types.Secret{secret})
	}

	secrets = append(secrets, secret)
	store.Set(SecretsStoreKey, secrets)
}

func SetDistro(distro types.Distro) {
	store.Set(DistroStoreKey, distro)
}

func SetParameters(params types.Parameters) {
	store.Set(ParametersStoreKey, params)
}

func SetImageInstance(image v1.Image) {
	store.Set(ImageInstanceStoreKey, image)
}

func SetScanElapsed(duration float64) {
	store.Set(ScanElapsedStoreKey, duration)
}
