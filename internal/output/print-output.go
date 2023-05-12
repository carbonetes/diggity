package output

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/internal/output/cyclonedx"
	"github.com/carbonetes/diggity/internal/output/github"
	"github.com/carbonetes/diggity/internal/output/save"
	"github.com/carbonetes/diggity/internal/output/spdx"
	"github.com/carbonetes/diggity/internal/output/tabular"
	"github.com/carbonetes/diggity/pkg/model"
	"golang.org/x/exp/maps"
)

var (
	// Result interface

	log = logger.GetLogger()
)

// PrintResults prints the result based on the arguments
func PrintResults(sbom *model.SBOM, args *model.Arguments) {
	Depulicate(sbom.Packages)
	SortResults(sbom.Packages)
	selectOutputType(args, sbom)
}

// Select Output Type based on the User Input with aliases considered
func selectOutputType(args *model.Arguments, results *model.SBOM) {
	for _, output := range strings.Split(strings.ToLower(args.Output.ToOutput()), ",") {
		switch output {
		case model.Table:
			tabular.PrintTable(args, results.Packages)
		case model.JSON.ToOutput():
			result, err := json.MarshalIndent(results, "", " ")
			if err != nil {
				log.Fatal(err)
			}
			if len(*args.OutputFile) > 0 {
				save.ResultToFile(string(result), args.OutputFile)
			} else {
				fmt.Printf("%+v\n", string(result))
			}
		case model.CycloneDXXML:
			cyclonedx.PrintCycloneDXXML(results.Packages, args.OutputFile)
		case model.CycloneDXJSON:
			cyclonedx.PrintCycloneDXJSON(results.Packages, args.OutputFile)
		case model.SPDXJSON:
			spdx.PrintSpdxJSON(args, results.Packages)
		case model.SPDXTagValue:
			spdx.PrintSpdxTagValue(args, results.Packages)
		case model.SPDXYML:
			spdx.PrintSpdxYaml(args, results.Packages)
		case model.GithubJSON:
			github.PrintGithubJSON(args, results)
		}
	}
}

// Remove Duplicates
func Depulicate(pkgs *[]model.Package) {
	result := make(map[string]model.Package, 0)
	for _, pkg := range *pkgs {
		if _, exists := result[pkg.Name+":"+pkg.Version+":"+pkg.Type]; !exists {
			result[pkg.Name+":"+pkg.Version+":"+pkg.Type] = pkg
		} else {
			idx := 0
			if len(pkg.Locations) > 0 {
				idx = len(pkg.Locations) - 1
				for _, l := range pkg.Locations {
					if l != pkg.Locations[idx] {
						pkg.Locations = append(pkg.Locations, model.Location{
							Path:      pkg.Path,
							LayerHash: "sha256:" + pkg.Locations[idx].LayerHash,
						})
						result[pkg.Name+":"+pkg.Version+":"+pkg.Type] = pkg
					}
				}
			}
		}
	}
	*pkgs = maps.Values(result)
}

func SortResults(pkgs *[]model.Package) {
	sort.SliceStable(*pkgs, func(i, j int) bool {
		if strings.EqualFold((*pkgs)[i].Name, (*pkgs)[j].Name) {
			return strings.ToLower((*pkgs)[i].Version) < strings.ToLower((*pkgs)[j].Version)
		}
		return strings.ToLower((*pkgs)[i].Name) < strings.ToLower((*pkgs)[j].Name)
	})
}
