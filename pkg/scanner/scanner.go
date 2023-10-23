// Package scan provides functionality for scanning Docker images.
package scanner

import (
	"log"

	"github.com/carbonetes/diggity/internal/slsa"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/common"
	diggity "github.com/carbonetes/diggity/pkg/parser"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

// Diggity scans the Docker images, Tar Files, and Codebases(directories) specified in the given model.Arguments struct and returns a sbom(model.Result) struct.
func Scan(arguments *model.Arguments) (*model.SBOM, *[]error) {
	params, err := common.NewParams(arguments)
	if err != nil {
		log.Fatal(err)
	}
	parsers := diggity.Parsers
	params.WG.Add(len(parsers))
	for _, parser := range parsers {
		parser(params)
	}

	params.WG.Wait()

	defer util.CleanUp(*params.DockerTemp)

	if *arguments.Provenance != "" {
		params.SBOM.SLSA = slsa.Provenance(params)
	}

	return params.SBOM, params.Errors
}
