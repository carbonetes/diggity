package nuget

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
	dotnetPackage = ".deps.json"
	dotnet        = "dotnet"
	nuget         = "nuget"
)

var dotnetMetadata metadata.DotnetDeps

// FindNugetPackagesFromContent - find nuget packages
func FindNugetPackagesFromContent(req *bom.ParserRequirements) {
	if util.ParserEnabled(nuget, req.Arguments.EnabledParsers) {
		for _, content := range *req.Contents {
			if strings.Contains(content.Path, dotnetPackage) {
				if err := parseNugetPackages(&content, req.SBOM.Packages); err != nil {
					err = errors.New("nuget-parser: " + err.Error())
					*req.Errors = append(*req.Errors, err)
				}
			}
		}
	}
	defer req.WG.Done()
}

// Parse nuget package metadata
func parseNugetPackages(location *model.Location, pkgs *[]model.Package) error {
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
				pkg := new(model.Package)
				pkg.ID = uuid.NewString()
				pkg.Name = split[0]
				pkg.Version = split[1]
				pkg.Type = dotnet
				pkg.Path = split[0]
				pkg.Locations = append(pkg.Locations, model.Location{
					Path:      util.TrimUntilLayer(*location),
					LayerHash: location.LayerHash,
				})
				parseNugetPURL(pkg)
				cpe.NewCPE23(pkg, pkg.Name, pkg.Name, pkg.Version)
				pkg.Metadata = cLib
				*pkgs = append(*pkgs, *pkg)
			}
		}
	}
	return nil
}

// Parse PURL
func parseNugetPURL(pkg *model.Package) {
	pkg.PURL = model.PURL("pkg" + ":" + "dotnet" + "/" + pkg.Name + "@" + pkg.Version)
}
