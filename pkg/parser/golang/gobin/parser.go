package gobin

import (
	"errors"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
)

const parserErr string = "gobin-parser: "

// Read go binaries content
func parseGoBinContent(location *model.Location, req *bom.ParserRequirements) {

	buildInfo, err := readFile(location.Path)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	if buildInfo == nil {
		return
	}

	if len(buildInfo.Main.Path) > 0 {
		pkg := newPackage(location, buildInfo, &buildInfo.Main)
		*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
	}

	if len(buildInfo.Deps) == 0 {
		return
	}

	// Parse dependencies
	for _, module := range buildInfo.Deps {
		if module.Replace != nil {
			pkg := newPackage(location, buildInfo, module.Replace)
			*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
		}
		pkg := newPackage(location, buildInfo, module)
		*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
	}
}
