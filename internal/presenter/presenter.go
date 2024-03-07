package presenter

import (
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/internal/presenter/json"
	"github.com/carbonetes/diggity/internal/presenter/status"
	"github.com/carbonetes/diggity/internal/presenter/table"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/types"
)

func DisplayResults(params types.Parameters, duration float64) {
	result := cdx.BOM

	format, saveToFile := params.OutputFormat, params.SaveToFile
	if !params.Quiet {
		status.Done()
	}
	switch format {
	case types.Table:
		table.Show(table.Create(), duration)
	case types.JSON, types.CycloneDXJSON, types.SPDXJSON:
		json.DisplayResults(result)
	case types.CycloneDXXML, types.SPDXXML:
		log.Error("XML output is not supported yet")
	default:
		log.Error("Unknown output format")
	}

	if len(saveToFile) > 0 {
		err := helper.SaveToFile(result, saveToFile, format.String())
		if err != nil {
			log.Errorf("Failed to save results to file : %s", err.Error())
		}
	}
}
