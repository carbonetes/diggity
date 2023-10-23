package composer

import (
	"errors"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/common"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const (
	composerLock string = "composer.lock"
	parserError  string = "cargo-parser: "
)

func parseComposerPackages(location *model.Location, req *common.ParserParams) {
	metadata, err := parseMetadata(location.Path)
	if err != nil {
		err = errors.New(parserError + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	if metadata == nil {
		return
	}

	if len(metadata.Packages) == 0 {
		return
	}
	for _, p := range metadata.Packages {
		pkg := newPackage(&p)
		if pkg == nil {
			continue
		}

		pkg.Metadata = p
		pkg.Locations = append(pkg.Locations, model.Location{
			Path:      util.TrimUntilLayer(*location),
			LayerHash: location.LayerHash,
		})
		*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
	}

	if len(metadata.PackagesDev) == 0 {
		return
	}

	for _, p := range metadata.PackagesDev {
		pkg := newPackage(&p)
		if pkg == nil {
			continue
		}
		pkg.Metadata = p
		pkg.Locations = append(pkg.Locations, model.Location{
			Path:      util.TrimUntilLayer(*location),
			LayerHash: location.LayerHash,
		})
		*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
	}
}
