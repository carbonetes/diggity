package parser

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/file"
	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/model/metadata"

	"github.com/google/uuid"
	"golang.org/x/mod/modfile"
)

const (
	goType   = "go"
	golang   = "golang"
	goModule = "go-module"
)

var (
	goModPath = filepath.Join("go.mod")
)

// GoModMetadata  metadata
type GoModMetadata map[string]interface{}

// FindGoModPackagesFromContent Find go.mod in the file contents
func FindGoModPackagesFromContent() {
	if parserEnabled(goType) {
		for _, content := range file.Contents {
			if filepath.Base(content.Path) == goModPath {
				if err := readGoModContent(content); err != nil {
					err = errors.New("go-mod-parser: " + err.Error())
					Errors = append(Errors, &err)
				}
			}
		}
	}
	defer WG.Done()
}

// Read go.mod content
func readGoModContent(location *model.Location) error {

	reader, err := os.Open(location.Path)
	if err != nil {
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
		_package := new(model.Package)
		_package = initGoModPackage(_package, location, modPkg)
		Packages = append(Packages, _package)
	}

	// Add New to Package List
	for _, modPkg := range goModFile.Replace {
		_package := new(model.Package)
		_package = initGoModPackage(_package, location, modPkg)
		Packages = append(Packages, _package)
	}

	// Cleanup Excluded Packages
	if len(goModFile.Exclude) > 0 {
		Packages = cleanExcluded(Packages, goModFile)
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
		Path:      TrimUntilLayer(*location),
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
func initGoModMetadata(_package *model.Package, modPkg interface{}) {
	var finalMetadata = metadata.GoModMetadata{}

	switch goPkg := modPkg.(type) {
	case *modfile.Require:
		finalMetadata.Path = goPkg.Mod.Path
		finalMetadata.Version = goPkg.Mod.Version
	case *modfile.Replace:
		finalMetadata.Path = goPkg.New.Path
		finalMetadata.Version = goPkg.New.Version
	}

	_package.Metadata = finalMetadata
}

// Parse PURL
func parseGoPackageURL(_package *model.Package) {
	_package.PURL = model.PURL(scheme + ":" + golang + "/" + _package.Name + "@" + _package.Version)
}

// Cleanup excluded packages
func cleanExcluded(packages []*model.Package, f *modfile.File) []*model.Package {
	for _, exPkg := range f.Exclude {
		for i, pkg := range packages {
			if pkg.Name == exPkg.Mod.Path {
				packages = append(packages[:i], packages[i+1:]...)
			}
		}
	}
	return packages
}

// Split go package path
func splitPath(path string) []string {
	return strings.Split(path, "/")
}
