package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/internal/output/cyclonedx"
	"github.com/carbonetes/diggity/internal/output/github"
	"github.com/carbonetes/diggity/internal/output/save"
	"github.com/carbonetes/diggity/internal/output/spdx"
	"github.com/carbonetes/diggity/internal/output/tabular"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"golang.org/x/exp/maps"
)

var (
	// Result interface

	log = logger.GetLogger()
)

// PrintResults prints the result based on the arguments
func PrintResults(req *bom.ParserRequirements) {
	Depulicate(req.Result.Packages)
	// SortResults(req.Result.Packages, result)
	// Table Output(Default)
	selectOutputType(req.Arguments, req.Result)

	if len(*req.Errors) > 0 {
		for _, err := range *req.Errors {
			log.Errorf("[warning]: %+v\n", err)
		}
	}
}

// Select Output Type based on the User Input with aliases considered
func selectOutputType(args *model.Arguments, results *model.Result) {
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

// Sort Results
// func SortResults(pkgs *[]model.Package, result map[string]model.Package) {
// 	*pkgs = maps.Values(result)
// 	sort.Slice(pkgs, func(i, j int) bool {
// 		if (*pkgs)[i].Name == (*pkgs)[j].Name {
// 			return (*pkgs)[i].Version < (*pkgs)[j].Version
// 		}
// 		return (*pkgs)[i].Name < (*pkgs)[j].Name
// 	})
// }
