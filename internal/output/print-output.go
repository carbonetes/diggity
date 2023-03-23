package output

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	log "github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/internal/output/cyclonedx"
	"github.com/carbonetes/diggity/internal/output/github"
	"github.com/carbonetes/diggity/internal/output/save"
	"github.com/carbonetes/diggity/internal/output/spdx"
	"github.com/carbonetes/diggity/internal/output/tabular"
	"github.com/carbonetes/diggity/internal/secret"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/distro"
	"github.com/carbonetes/diggity/pkg/parser/docker"

	"golang.org/x/exp/maps"
)

type result map[string]*model.Package

// Result interface
var Result result = make(map[string]*model.Package, 0)

// PrintResults prints the result based on the arguments
func PrintResults() {
	finalizeResults()
	outputTypes := strings.ToLower(bom.Arguments.Output.ToOutput())

	// Table Output(Default)
	selectOutputType(outputTypes)

	if len(bom.Errors) > 0 {
		for _, err := range bom.Errors {
			log.GetLogger().Printf("[warning]: %+v\n", *err)
		}
	}
}

// Select Output Type based on the User Input with aliases considered
func selectOutputType(outputTypes string) {
	for _, output := range strings.Split(outputTypes, ",") {
		switch output {
		case model.Table:
			tabular.PrintTable()
		case model.JSON.ToOutput():
			if len(*bom.Arguments.OutputFile) > 0 {
				save.ResultToFile(GetResults())
			} else {
				fmt.Printf("%+v\n", GetResults())
			}
		case model.CycloneDXXML:
			cyclonedx.PrintCycloneDXXML()
		case model.CycloneDXJSON:
			cyclonedx.PrintCycloneDXJSON()
		case model.SPDXJSON:
			spdx.PrintSpdxJSON()
		case model.SPDXTagValue:
			spdx.PrintSpdxTagValue()
		case model.GithubJSON:
			github.PrintGithubJSON()
		}
	}
}

// Remove Duplicates and Sort Results
func finalizeResults() {
	for _, pkg := range bom.Packages {
		if _, exists := Result[pkg.Name+":"+pkg.Version+":"+pkg.Type]; !exists {
			Result[pkg.Name+":"+pkg.Version+":"+pkg.Type] = pkg
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
						Result[pkg.Name+":"+pkg.Version+":"+pkg.Type] = pkg
					}
				}
			}
		}
	}
	sortResults()
}

// Sort Results
func sortResults() {
	bom.Packages = maps.Values(Result)
	sort.Slice(bom.Packages, func(i, j int) bool {
		if bom.Packages[i].Name == bom.Packages[j].Name {
			return bom.Packages[i].Version < bom.Packages[j].Version
		}
		return bom.Packages[i].Name < bom.Packages[j].Name
	})
}

// GetResults - For event bus handler
func GetResults() string {
	pkgs := maps.Values(Result)

	sort.Slice(pkgs, func(i, j int) bool {
		return pkgs[i].Name < pkgs[j].Name
	})

	output := model.Result{
		Distro:   distro.Distro(),
		Packages: bom.Packages,
	}

	if !*bom.Arguments.DisableSecretSearch {
		output.Secret = secret.SecretResults
	}

	output.ImageInfo = docker.ImageInfo

	result, _ := json.MarshalIndent(output, "", " ")
	return string(result)
}
