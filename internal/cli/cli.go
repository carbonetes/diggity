package cli

import (
	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/internal/output"
	"github.com/carbonetes/diggity/internal/ui"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scanner"
)

var log = logger.GetLogger()

func Start(arguments *model.Arguments) {
	sbom, errs := scanner.Scan(arguments)
	ui.DoneSpinner()
	if errs != nil {
		if len(*errs) > 0 {
			for _, err := range *errs {
				log.Warningln(err.Error())
			}
		}
	}
	output.PrintResults(sbom, arguments)
}
