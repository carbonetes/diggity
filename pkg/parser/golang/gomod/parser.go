package gomod

import (
	"errors"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"golang.org/x/mod/modfile"
)

const parserErr = "gomod-parser: "

// Read go.mod content
func parseGoModContent(location *model.Location, req *bom.ParserRequirements) {

	modFile, err := readModFile(location.Path)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	if modFile == nil {
		return
	}

	// Initial Package List
	for _, modPkg := range modFile.Require {
		pkg := newPackage(location, modPkg)
		*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
	}

	// Add New to Package List
	for _, modPkg := range modFile.Replace {
		pkg := newPackage(location, modPkg)
		*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
	}

	// Cleanup Excluded Packages
	if len(modFile.Exclude) > 0 {
		cleanExcluded(req.SBOM.Packages, modFile)
	}
}

// Cleanup excluded packages
func cleanExcluded(packages *[]model.Package, f *modfile.File) {
	for _, exPkg := range f.Exclude {
		for i, pkg := range *packages {
			if pkg.Name != exPkg.Mod.Path {
				continue
			}
			*packages = append((*packages)[:i], (*packages)[i+1:]...)
		}
	}
}
