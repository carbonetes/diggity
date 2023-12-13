package stream

import (
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/types"
)

func AggrerateSoftwareManifest() types.SoftwareManifest {
	params := GetParameters()
	var sbom interface{}
	var files []string
	switch params.OutputFormat {
	case types.CycloneDXJSON, types.CycloneDXXML:
	case types.SPDXJSON, types.SPDXTag, types.SPDXXML:
	default:
		sbom = AggregateSBOM()
	}
	if params.AllowFileListing {
		files = GetFiles()
	}

	return types.SoftwareManifest{
		SBOM:       sbom,
		ImageInfo:  GetImageInfo(),
		Distro:     GetDistro(),
		Secrets:    AggregateSecrets(),
		Files:      files,
		Parameters: params,
	}
}

func AggregateSBOM() types.SBOM {
	data, _ := store.Get(SBOMStoreKey)

	sbom, ok := data.(types.SBOM)

	if !ok {
		log.Error("AggregateSBOM received unknown data type")
	}

	sbom.Components = append(sbom.Components, GetComponents()...)
	sbom.Total = len(sbom.Components)
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
