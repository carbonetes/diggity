package stream

import (
	"github.com/carbonetes/diggity/internal/log"
	convert "github.com/carbonetes/diggity/pkg/converter/cdx"
	"github.com/carbonetes/diggity/pkg/converter/spdx"
	"github.com/carbonetes/diggity/pkg/provenance"
	"github.com/carbonetes/diggity/pkg/types"
)

func AggrerateSoftwareManifest() types.SoftwareManifest {
	params := GetParameters()
	var sbom interface{}
	var files []string
	var slsa types.SLSA
	switch params.OutputFormat {
	case types.CycloneDXJSON, types.CycloneDXXML:
		bom := AggregateSBOM()
		cdx := convert.ToCDX(&bom)
		sbom = cdx

	case types.SPDXJSON, types.SPDXTag, types.SPDXXML:
		bom := AggregateSBOM()
		sbom = spdx.ToSPDX23(bom, params.Input)
	default:
		sbom = AggregateSBOM()
	}
	if params.AllowFileListing {
		files = GetFiles()
	}

	if len(params.Provenance) != 0 {
		p, err := provenance.Parse(params.Provenance)
		if err != nil {
			log.Error(err)
		} else {
			slsa.Provenance = p
		}
	}

	return types.SoftwareManifest{
		SBOM:       sbom,
		ImageInfo:  GetImageInfo(),
		Distro:     GetDistro(),
		Secrets:    AggregateSecrets(),
		Files:      files,
		Parameters: params,
		SLSA:       slsa,
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
