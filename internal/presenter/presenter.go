package presenter

import (
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/internal/presenter/json"
	"github.com/carbonetes/diggity/internal/presenter/status"
	"github.com/carbonetes/diggity/internal/presenter/table"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/types"
	"github.com/golistic/urn"
)

func DisplayResults(params types.Parameters, duration float64, addr *urn.URN) {
	result := cdx.Finalize(addr)
	format, filename := params.OutputFormat, params.SaveToFile
	if !params.Quiet {
		status.Done()
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
		table.Show(table.Create(result), duration)
	case types.JSON, types.CycloneDXJSON, types.SPDXJSON:
		json.DisplayResults(result)
	case types.CycloneDXXML, types.SPDXXML:
		log.Debug("XML output is not supported yet")
	default:
		log.Debug("Unknown output format")
	}

}
