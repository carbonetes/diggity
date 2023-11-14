package stream

import (
	"log"

	"github.com/carbonetes/diggity/pkg/types"
)

func AggrerateSoftwareManifest() types.SoftwareManifest {
	return types.SoftwareManifest{
		SBOM:          AggregateSBOM(),
		ImageManifest: GetImageManifest(),
		Distro:        GetDistro(),
		Secret:        AggregateSecrets(),
		Parameters:    GetParameters(),
	}
}

func AggregateSBOM() types.SBOM {
	data, _ := store.Get(SBOMStoreKey)

	sbom, ok := data.(types.SBOM)

	if !ok {
		log.Fatal("AggregateSBOM received unknown data type")
	}

	sbom.Components = append(sbom.Components, GetComponents()...)

	store.Set(SBOMStoreKey, sbom)
	return sbom
}

func AggregateSecrets() types.SecretResult {
	data, _ := store.Get(SecretsStoreKey)

	secrets, ok := data.([]types.Secret)

	if !ok {
		log.Fatal("AggregateSecrets received unknown data type")
	}

	return types.SecretResult{
		Parameters: GetSecretParameters(),
		Secrets:    secrets,
	}

}
