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
	dotnetPackage = ".deps.json"
	dotnet        = "dotnet"
	nuget         = "nuget"
)

var dotnetMetadata metadata.DotnetDeps

// FindNugetPackagesFromContent - find nuget packages
func FindNugetPackagesFromContent() {
	if parserEnabled(nuget) {
		for _, content := range file.Contents {
			if strings.Contains(content.Path, dotnetPackage) {
				if err := parseNugetPackages(content); err != nil {
					err = errors.New("nuget-parser: " + err.Error())
					Errors = append(Errors, &err)
				}
			}
		}
	}
	defer WG.Done()
}

// Parse nuget package metadata
func parseNugetPackages(location *model.Location) error {
	byteValue, err := os.ReadFile(location.Path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(byteValue, &dotnetMetadata); err != nil {
		return err
	}
	if len(dotnetMetadata.Libraries) > 0 {

		for nameAndVersion, cLib := range dotnetMetadata.Libraries {

			if cLib.Type == "package" {
				split := strings.Split(nameAndVersion, "/")
				_package := new(model.Package)
				_package.ID = uuid.NewString()
				_package.Name = split[0]
				_package.Version = split[1]
				_package.Type = dotnet
				_package.Path = split[0]
				_package.Locations = append(_package.Locations, model.Location{
					Path:      TrimUntilLayer(*location),
					LayerHash: location.LayerHash,
				})
				parseNugetPURL(_package)
				cpe.NewCPE23(_package, _package.Name, _package.Name, _package.Version)
				_package.Metadata = cLib
				Packages = append(Packages, _package)
			}
		}
	}
	return nil
}

// Parse PURL
func parseNugetPURL(_package *model.Package) {
	_package.PURL = model.PURL(scheme + ":" + "dotnet" + "/" + _package.Name + "@" + _package.Version)
}
