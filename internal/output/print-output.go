package output

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	log "github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/parser"
	"github.com/carbonetes/diggity/internal/secret"

	"golang.org/x/exp/maps"
)

type result map[string]*model.Package

// Result interface
var Result result = make(map[string]*model.Package, 0)

// PrintResults prints the result based on the arguments
func PrintResults() {
	finalizeResults()
	outputTypes := strings.ToLower(parser.Arguments.Output.ToOutput())

	// Table Output(Default)
	selectOutputType(outputTypes)

	if len(parser.Errors) > 0 {
		for _, err := range parser.Errors {
			log.GetLogger().Printf("[warning]: %+v\n", *err)
		}
	}
}

// Select Output Type based on the User Input with aliases considered
func selectOutputType(outputTypes string) {
	for _, output := range strings.Split(outputTypes, ",") {
		switch output {
		case model.Table:
			printTable()
		case model.JSON.ToOutput():
			if len(*parser.Arguments.OutputFile) > 0 {
				saveResultToFile(GetResults())
			} else {
				fmt.Printf("%+v\n", GetResults())
			}
		case model.CycloneDXXML, "cyclonedxxml", "cyclonedx", "cyclone":
			printCycloneDXXML()
		case model.CycloneDXJSON, "cyclonedxjson":
			printCycloneDXJSON()
		case model.SPDXJSON, "spdxjson":
			printSpdxJSON()
		case model.SPDXTagValue, "spdxtagvalue", "spdx", "spdxtv":
			printSpdxTagValue()
		}
	}
}

// Remove Duplicates and Sort Results
func finalizeResults() {
	for _, _package := range parser.Packages {
		if _, exists := Result[_package.Name+":"+_package.Version+":"+_package.Type]; !exists {
			Result[_package.Name+":"+_package.Version+":"+_package.Type] = _package
		} else {
			idx := 0
			if len(_package.Locations) > 0 {
				idx = len(_package.Locations) - 1
				for _, l := range _package.Locations {
					if l != _package.Locations[idx] {
						_package.Locations = append(_package.Locations, model.Location{
							Path:      _package.Path,
							LayerHash: "sha256:" + _package.Locations[idx].LayerHash,
						})
						Result[_package.Name+":"+_package.Version+":"+_package.Type] = _package
					}
				}
			}
		}
	}
	sortResults()
}

// Sort Results
func sortResults() {
	parser.Packages = maps.Values(Result)
	sort.Slice(parser.Packages, func(i, j int) bool {
		if parser.Packages[i].Name == parser.Packages[j].Name {
			return parser.Packages[i].Version < parser.Packages[j].Version
		}
		return parser.Packages[i].Name < parser.Packages[j].Name
	})
}

// GetResults - For event bus handler
func GetResults() string {
	_packages := maps.Values(Result)

	sort.Slice(_packages, func(i, j int) bool {
		return _packages[i].Name < _packages[j].Name
	})

	output := Output{
		Distro:   parser.Distro(),
		Packages: parser.Packages,
	}

	if !*parser.Arguments.DisableSecretSearch {
		output.Secret = secret.SecretResults
	}

	output.ImageInfo = parser.ImageInfo

	result, _ := json.MarshalIndent(output, "", " ")
	return string(result)
}
