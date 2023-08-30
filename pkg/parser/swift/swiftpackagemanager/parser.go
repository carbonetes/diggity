package swiftpackagemanager

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

func parseSwiftPackages(location *model.Location, req *bom.ParserRequirements) {
	file, err := os.Open(location.Path)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
		return
	}

	defer file.Close()

	var packagemetadata metadata.SwiftPackageManagerMetadata
	err = json.NewDecoder(file).Decode(&packagemetadata)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
		return
	}

	switch packagemetadata.Version {
	case 1:
		for _, pin := range packagemetadata.Object.Pins {
			pkg := newV1Package(&pin)
			pkg.Path = util.TrimUntilLayer(*location)
			pkg.Locations = append(pkg.Locations, model.Location{
				Path:      pkg.Path,
				LayerHash: location.LayerHash,
			})
			*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
		}
	case 2:
		for _, pin := range packagemetadata.Pins {
			pkg := newV2Package(&pin)
			pkg.Path = util.TrimUntilLayer(*location)
			pkg.Locations = append(pkg.Locations, model.Location{
				Path:      pkg.Path,
				LayerHash: location.LayerHash,
			})
			*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
		}

	}

}
