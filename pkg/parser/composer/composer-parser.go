package composer

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"

	"github.com/google/uuid"
)

const (
	phpType      = "php"
	composerLock = "composer.lock"
	composer     = "composer"
)

var lockFileMetadata metadata.ComposerMetadata

// FindComposerPackagesFromContent - find composers packages from content
func FindComposerPackagesFromContent(req *bom.ParserRequirements) {
	if util.ParserEnabled(composer, req.Arguments.EnabledParsers) {
		for _, content := range *req.Contents {
			if strings.Contains(content.Path, composerLock) {
				if err := parseComposerPackages(&content, req.Result.Packages); err != nil {
					err = errors.New("composer-parser: " + err.Error())
					*req.Errors = append(*req.Errors, err)
				}
			}
		}
	}
	defer req.WG.Done()
}

// Parse composer package metadata
func parseComposerPackages(location *model.Location, pkgs *[]model.Package) error {
	byteValue, err := os.ReadFile(location.Path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(byteValue, &lockFileMetadata); err != nil {
		return err
	}

	for _, cPackage := range lockFileMetadata.Packages {
		pkg := new(model.Package)
		pkg.ID = uuid.NewString()
		pkg.Name = cPackage.Name
		pkg.Version = cPackage.Version
		pkg.Description = cPackage.Description
		pkg.Licenses = cPackage.License
		pkg.Type = phpType
		pkg.Path = cPackage.Name
		pkg.Metadata = cPackage
		pkg.Locations = append(pkg.Locations, model.Location{
			Path:      util.TrimUntilLayer(*location),
			LayerHash: location.LayerHash,
		})
		parseComposerPURL(pkg)
		vendorProduct := strings.Split(pkg.Name, "/")
		if len(vendorProduct) == 0 {
			vendorProduct = []string{
				pkg.Name,
				pkg.Name,
			}
		}
		cpe.NewCPE23(pkg, vendorProduct[0], vendorProduct[1], pkg.Version)
		*pkgs = append(*pkgs, *pkg)
	}

	for _, cPackage := range lockFileMetadata.PackagesDev {
		pkg := new(model.Package)
		pkg.ID = uuid.NewString()
		pkg.Name = cPackage.Name
		pkg.Version = cPackage.Version
		pkg.Description = cPackage.Description
		pkg.Licenses = cPackage.License
		pkg.Type = phpType
		pkg.Path = cPackage.Name
		pkg.Metadata = cPackage
		pkg.Locations = append(pkg.Locations, model.Location{
			Path:      util.TrimUntilLayer(*location),
			LayerHash: location.LayerHash,
		})
		parseComposerPURL(pkg)
		vendorProduct := strings.Split(pkg.Name, "/")
		if len(vendorProduct) == 0 {
			vendorProduct = []string{
				pkg.Name,
				pkg.Name,
			}
		}
		cpe.NewCPE23(pkg, vendorProduct[0], vendorProduct[1], pkg.Version)
		*pkgs = append(*pkgs, *pkg)
	}

	return nil
}

// Parse PURL
func parseComposerPURL(pkg *model.Package) {
	pkg.PURL = model.PURL("pkg" + ":" + "composer" + "/" + pkg.Name + "@" + pkg.Version)
}
