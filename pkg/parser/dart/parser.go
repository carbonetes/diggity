package dart

import (
	"errors"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

// Parse dart package metadata
func parseDartPackages(location *model.Location, req *bom.ParserRequirements) {
	metadata, err := parseMetadata(location.Path)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	if metadata == nil {
		return
	}

	pkg := newPackage(*metadata)

	if pkg == nil {
		return
	}

	if len(pkg.Name) == 0 {
		return
	}

	pkg.Locations = append(pkg.Locations, model.Location{
		Path:      util.TrimUntilLayer(*location),
		LayerHash: location.LayerHash,
	})

	*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
}

// Parse dart packages metadata - lock file
func parseDartLockPackages(location *model.Location, req *bom.ParserRequirements) {
	metadata, err := parseLockfileMetadata(location.Path)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	if metadata == nil {
		return
	}

	if len(metadata.Packages) == 0 {
		return
	}

	for _, dartMetadata := range metadata.Packages {
		pkg := newPackage(dartMetadata)

		if pkg == nil {
			continue
		}

		if len(pkg.Name) == 0 {
			continue
		}

		pkg.Locations = append(pkg.Locations, model.Location{
			Path:      util.TrimUntilLayer(*location),
			LayerHash: location.LayerHash,
		})

		*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
	}
}
