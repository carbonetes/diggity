package stream

import (
	"github.com/carbonetes/diggity/pkg/types"
)

func AggrerateSoftwareManifest() types.SoftwareManifest {
	return types.SoftwareManifest{
		SBOM:       AggregateSBOM(),
		ImageInfo:  GetImageInfo(),
		Distro:     GetDistro(),
		Secrets:    AggregateSecrets(),
		Parameters: GetParameters(),
	}
}

func AggregateSBOM() types.SBOM {
	data, _ := store.Get(SBOMStoreKey)

	sbom, ok := data.(types.SBOM)

	if !ok {
		log.Error("AggregateSBOM received unknown data type")
	}

	sbom.Components = append(sbom.Components, GetComponents()...)

	store.Set(SBOMStoreKey, sbom)
	return sbom
}

func AggregateSecrets() []types.Secret {
	data, _ := store.Get(SecretsStoreKey)

	secrets, ok := data.([]types.Secret)

	if !ok {
		log.Error("AggregateSecrets received unknown data type")
	}

	return secrets

}
