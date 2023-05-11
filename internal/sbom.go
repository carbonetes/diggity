package sbom

import (
	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/internal/output"
	"github.com/carbonetes/diggity/internal/slsa"
	"github.com/carbonetes/diggity/internal/ui"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

var log = logger.GetLogger()

// Start SBOM extraction
func Start(arguments *model.Arguments) {
	if *arguments.Quiet {
		ui.Disable()
	}
	pb := ui.InitSpinner("Scanning for packages...")
	go ui.RunSpinner(pb)

	requirements, err := bom.InitParsers(arguments)
	if err != nil {
		log.Fatal(err)
	}
	requirements.WG.Add(len(parser.FindFunctions))
	for _, parser := range parser.FindFunctions {
		go parser(requirements)
	}
	requirements.WG.Wait()
	defer util.CleanUp(requirements)

	if *arguments.Provenance != "" {
		requirements.SBOM.SLSA = slsa.Provenance(requirements)
	}

	ui.DoneSpinner(pb)

	//Print Results and Cleanup
	output.PrintResults(requirements)
}
