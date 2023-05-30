package rpm

import (
	"errors"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	rpmdb "github.com/knqyf263/go-rpmdb/pkg"
)

// Read RPM package information from rpm/Packages
func readRpmContent(location *model.Location, req *bom.ParserRequirements) {

	// Open and Get rpm/Packages data
	db, err := rpmdb.Open(location.Path)
	if err != nil {
		err = errors.New("rpm-parser: " + err.Error())
		*req.Errors = append(*req.Errors, err)
	}
	rpmPkgList, err := db.ListPackages()
	if err != nil {
		err = errors.New("rpm-parser: " + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	for _, rpmPkg := range rpmPkgList {
		// Get RPM package contents
		pkg := new(model.Package)
		pkg = initRpmPackage(pkg, location, rpmPkg)

		*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
	}
}
