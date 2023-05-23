package nuget

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const parserErr string = "nuget-parser: "

// Parse nuget package metadata
func parseNugetPackages(location *model.Location, req *bom.ParserRequirements) {
	byteValue, err := os.ReadFile(location.Path)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	var metadata metadata.DotnetDeps
	if err := json.Unmarshal(byteValue, &metadata); err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	if len(metadata.Libraries) > 0 {

		for keyValue, lib := range metadata.Libraries {
			if lib.Type != "package" {
				continue
			}

			pkg := newPackage(keyValue, &lib)

			if pkg == nil {
				continue
			}

			pkg.Locations = append(pkg.Locations, model.Location{
				Path:      util.TrimUntilLayer(*location),
				LayerHash: location.LayerHash,
			})

			*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
		}
	}
}
