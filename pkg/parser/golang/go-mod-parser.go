package golang

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"

	"github.com/google/uuid"
	"golang.org/x/mod/modfile"
)

const (
	goType       = "go"
	golang       = "golang"
	goModule     = "go-module"
	noFileErrWin = "The system cannot find the file specified"
	noFileErrMac = "no such file or directory"
)

var (
	goModPath = filepath.Join("go.mod")
)

// GoModMetadata  metadata
type GoModMetadata map[string]interface{}

// FindGoModPackagesFromContent Find go.mod in the file contents
func FindGoModPackagesFromContent(req *bom.ParserRequirements) {
	if util.ParserEnabled(goType, req.Arguments.EnabledParsers) {
		for _, content := range *req.Contents {
			if filepath.Base(content.Path) == goModPath {
				if err := readGoModContent(&content, req.SBOM.Packages); err != nil {
					err = errors.New("go-mod-parser: " + err.Error())
					*req.Errors = append(*req.Errors, err)
				}
			}
		}
	}
	defer req.WG.Done()
}

// Read go.mod content
func readGoModContent(location *model.Location, pkgs *[]model.Package) error {

	reader, err := os.Open(location.Path)
	if err != nil {
		if strings.Contains(err.Error(), noFileErrWin) || strings.Contains(err.Error(), noFileErrMac) {
			return nil
		}
		return err
	}
	defer reader.Close()

	modContents, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	goModFile, err := modfile.Parse(location.Path, modContents, nil)
	if err != nil {
		return err
	}

	// Initial Package List
	for _, modPkg := range goModFile.Require {
		pkg := new(model.Package)
		pkg = initGoModPackage(pkg, location, modPkg)
		*pkgs = append(*pkgs, *pkg)
	}

	// Add New to Package List
	for _, modPkg := range goModFile.Replace {
		pkg := new(model.Package)
		pkg = initGoModPackage(pkg, location, modPkg)
		*pkgs = append(*pkgs, *pkg)
	}

	// Cleanup Excluded Packages
	if len(goModFile.Exclude) > 0 {
		cleanExcluded(pkgs, goModFile)
	}

	return nil
}

// Initialize go mod package contents
func initGoModPackage(p *model.Package, location *model.Location, modPkg interface{}) *model.Package {
	p.ID = uuid.NewString()

	switch goPkg := modPkg.(type) {
	case *modfile.Require:
		p.Name = goPkg.Mod.Path
		p.Version = goPkg.Mod.Version
	case *modfile.Replace:
		p.Name = goPkg.New.Path
		p.Version = goPkg.New.Version
	}

	p.Type = goModule
	p.Path = p.Name
	p.Licenses = []string{}

	// get locations
	p.Locations = append(p.Locations, model.Location{
		Path:      util.TrimUntilLayer(*location),
		LayerHash: location.LayerHash,
	})

	// get CPEs
	cpePaths := splitPath(p.Name)

	// check if cpePaths only contains the product
	if len(cpePaths) > 1 {
		cpe.NewCPE23(p, cpePaths[len(cpePaths)-2], cpePaths[len(cpePaths)-1], p.Version)
	} else {
		cpe.NewCPE23(p, "", cpePaths[0], p.Version)
	}

	// get purl
	parseGoPackageURL(p)

	// parse and fill final metadata
	initGoModMetadata(p, modPkg)

	return p
}

// Initialize go metadata values from content
func initGoModMetadata(pkg *model.Package, modPkg interface{}) {
	var finalMetadata = metadata.GoModMetadata{}

	switch goPkg := modPkg.(type) {
	case *modfile.Require:
		finalMetadata.Path = goPkg.Mod.Path
		finalMetadata.Version = goPkg.Mod.Version
	case *modfile.Replace:
		finalMetadata.Path = goPkg.New.Path
		finalMetadata.Version = goPkg.New.Version
	}

	pkg.Metadata = finalMetadata
}

// Parse PURL
func parseGoPackageURL(pkg *model.Package) {
	pkg.PURL = model.PURL("pkg" + ":" + golang + "/" + pkg.Name + "@" + pkg.Version)
}

// Cleanup excluded packages
func cleanExcluded(packages *[]model.Package, f *modfile.File) {
	for _, exPkg := range f.Exclude {
		for i, pkg := range *packages {
			if pkg.Name == exPkg.Mod.Path {
				*packages = append((*packages)[:i], (*packages)[i+1:]...)
			}
		}
	}
}

// Split go package path
func splitPath(path string) []string {
	return strings.Split(path, "/")
}
