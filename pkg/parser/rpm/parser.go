package rpm

import (
	"errors"
	"io"
	"os"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
	rpmdb "github.com/knqyf263/go-rpmdb/pkg"
	_ "modernc.org/sqlite"
)

// Read RPM package information from rpm/Packages
func readRpmContent(location *model.Location, req *bom.ParserRequirements) {
	f, err := os.Open(location.Path)
	if err != nil {
		err = errors.New("rpm-parser: " + err.Error())
		*req.Errors = append(*req.Errors, err)
		return
	}

	defer f.Close()

	tmp, err := os.CreateTemp("", "rpmdb_")
	if err != nil {
		err = errors.New("rpm-parser: " + err.Error())
		*req.Errors = append(*req.Errors, err)
		return
	}

	defer util.CleanUp(tmp.Name())

	_, err = io.Copy(tmp, f)
	if err != nil {
		err = errors.New("rpm-parser: " + err.Error())
		*req.Errors = append(*req.Errors, err)
		return
	}

	defer tmp.Close()

	db, err := rpmdb.Open(tmp.Name())

	if err != nil {
		err = errors.New("rpm-parser: " + err.Error())
		*req.Errors = append(*req.Errors, err)
		return
	}

	rpmPkgList, err := db.ListPackages()
	if err != nil {
		err = errors.New("rpm-parser: " + err.Error())
		*req.Errors = append(*req.Errors, err)
		return
	}

	if len(rpmPkgList) == 0 {
		return
	}

	for _, pkgInfo := range rpmPkgList {
		if pkgInfo == nil {
			continue
		}
		// Get RPM package contents
		pkg := initRpmPackage(location, pkgInfo)
		pkg.Path = util.TrimUntilLayer(*location)
		// get locations
		pkg.Locations = append(pkg.Locations, model.Location{
			Path:      pkg.Path,
			LayerHash: location.LayerHash,
		})

		*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
	}
}
