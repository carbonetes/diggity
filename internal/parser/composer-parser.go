package parser

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/file"
	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/model/metadata"

	"github.com/google/uuid"
)

const (
	phpType      = "php"
	composerLock = "composer.lock"
	composer     = "composer"
)

var lockFileMetadata metadata.ComposerMetadata

// FindComposerPackagesFromContent - find composers packages from content
func FindComposerPackagesFromContent() {
	if parserEnabled(composer) {
		for _, content := range file.Contents {
			if strings.Contains(content.Path, composerLock) {
				if err := parseComposerPackages(content); err != nil {
					err = errors.New("composer-parser: " + err.Error())
					Errors = append(Errors, &err)
				}
			}
		}
	}
	defer WG.Done()
}

// Parse composer package metadata
func parseComposerPackages(location *model.Location) error {
	byteValue, err := os.ReadFile(location.Path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(byteValue, &lockFileMetadata); err != nil {
		return err
	}

	for _, cPackage := range lockFileMetadata.Packages {
		_package := new(model.Package)
		_package.ID = uuid.NewString()
		_package.Name = cPackage.Name
		_package.Version = cPackage.Version
		_package.Description = cPackage.Description
		_package.Licenses = cPackage.License
		_package.Type = phpType
		_package.Path = cPackage.Name
		_package.Metadata = cPackage
		_package.Locations = append(_package.Locations, model.Location{
			Path:      TrimUntilLayer(*location),
			LayerHash: location.LayerHash,
		})
		parseComposerPURL(_package)
		vendorProduct := strings.Split(_package.Name, "/")
		if len(vendorProduct) == 0 {
			vendorProduct = []string{
				_package.Name,
				_package.Name,
			}
		}
		cpe.NewCPE23(_package, vendorProduct[0], vendorProduct[1], _package.Version)
		Packages = append(Packages, _package)
	}

	for _, cPackage := range lockFileMetadata.PackagesDev {
		_package := new(model.Package)
		_package.ID = uuid.NewString()
		_package.Name = cPackage.Name
		_package.Version = cPackage.Version
		_package.Description = cPackage.Description
		_package.Licenses = cPackage.License
		_package.Type = phpType
		_package.Path = cPackage.Name
		_package.Metadata = cPackage
		_package.Locations = append(_package.Locations, model.Location{
			Path:      TrimUntilLayer(*location),
			LayerHash: location.LayerHash,
		})
		parseComposerPURL(_package)
		vendorProduct := strings.Split(_package.Name, "/")
		if len(vendorProduct) == 0 {
			vendorProduct = []string{
				_package.Name,
				_package.Name,
			}
		}
		cpe.NewCPE23(_package, vendorProduct[0], vendorProduct[1], _package.Version)
		Packages = append(Packages, _package)
	}

	return nil
}

// Parse PURL
func parseComposerPURL(_package *model.Package) {
	_package.PURL = model.PURL(scheme + ":" + "composer" + "/" + _package.Name + "@" + _package.Version)
}
