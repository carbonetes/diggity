package cocoapods

import (
	"errors"
	"os"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
	"gopkg.in/yaml.v3"
)

const parserErr string = "cocoapods-parser: "

func parseSwiftPackages(location *model.Location, req *bom.ParserRequirements) {
	byteValue, err := os.ReadFile(location.Path)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}
	var podFileLockFileMetadata metadata.PodFileLockMetadata
	if err := yaml.Unmarshal([]byte(byteValue), &podFileLockFileMetadata); err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	for _, pod := range podFileLockFileMetadata.Pods {
		pkg := newPackage(pod, &podFileLockFileMetadata)
		if pkg == nil {
			continue
		}
		if pkg.Name == "" {
			continue
		}

		pkg.Locations = append(pkg.Locations, model.Location{
			Path:      util.TrimUntilLayer(*location),
			LayerHash: location.LayerHash,
		})
		*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
	}
}
