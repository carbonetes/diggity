package ui

import (
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/types"
	"github.com/golistic/urn"
)

func DisplayResult(params types.Parameters, duration float64, addr *urn.URN) {
	result := cdx.Finalize(addr)
	format, filename := params.OutputFormat, params.SaveToFile
	if !params.Quiet {
		Done()
	}

	if len(filename) > 0 {
		err := helper.SaveToFile(result, filename, string(types.JSON))
		if err != nil {
			log.Debugf("Failed to save results to file : %s", err.Error())
		}
		return
	}

	switch format {
	case types.Table:
		ShowTable(CreateTable(result), duration)
	case types.JSON, types.CycloneDXJSON, types.SPDXJSON:
		printJson(result)
	case types.CycloneDXXML, types.SPDXXML:
		log.Debug("XML output is not supported yet")
	default:
		log.Debug("Unknown output format")
	}

}
