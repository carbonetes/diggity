package presenter

import (
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/internal/presenter/json"
	"github.com/carbonetes/diggity/internal/presenter/table"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

func DisplayResults(data interface{}) interface{} {
	duration, ok := data.(float64)
	if !ok {
		log.Error("DisplayResults received unknown type")
	}

	params := stream.GetParameters()
	result := stream.AggrerateSoftwareManifest()
	result.Duration = duration
	format, saveToFile := params.OutputFormat, params.SaveToFile

	if len(saveToFile) > 0 {
		err := helper.SaveToFile(result, saveToFile, format.String())
		if err != nil {
			log.Errorf("Failed to save results to file : %s", err.Error())
		}
		return data
	}

	switch format {
	case types.Table:
		table.Show(table.Create(), duration)
	case types.JSON:
		json.DisplayResults(result)
	default:
		log.Error("Unknown output format")
	}

	return data
}
