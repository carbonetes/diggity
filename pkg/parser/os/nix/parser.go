package nix

import (
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/common"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

func parseNixPackage(input string, location *model.Location, req *common.ParserParams) {
	metadata := parseNixPath(input)
	if metadata == nil {
		return
	}

	pkg := newPackage(metadata)

	if len(pkg.Name) == 0 || len(pkg.Version) == 0 {
		return
	}
	generateCpes(pkg)
	pkg.Path = util.TrimUntilLayer(*location)
	pkg.Locations = append(pkg.Locations, model.Location{
		Path:      pkg.Path,
		LayerHash: location.LayerHash,
	})

	*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
}
