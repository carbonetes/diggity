package portage

import (
	"errors"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/common"
)

var (
	portageContent   = "CONTENTS"
	portageLicense   = "LICENSE"
	portageSize      = "SIZE"
	portageObj       = "obj"
	portageAlgorithm = "md5"
	ebuild           = "ebuild"
	noFileErrWin     = "The system cannot find the file specified"
	noFileErrMac     = "no such file or directory"
)

// Read Portage Contents
func readPortageContent(location *model.Location, req *common.ParserParams) {
	// Parse package metadata from path
	pkg, err := initPortagePackage(location, req.Arguments.DisableFileListing)
	if err != nil {
		err = errors.New("portage-parser: " + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	if pkg == nil {
		return
	}

	if pkg.Name == "" {
		return
	}

	*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
}
