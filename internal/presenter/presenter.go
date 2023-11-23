package presenter

import (
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/internal/presenter/json"
	"github.com/carbonetes/diggity/internal/presenter/table"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

var log = logger.GetLogger()

func DisplayResults(data interface{}) interface{} {
	duration, ok := data.(float64)
	if !ok {
		log.Fatal("DisplayResults received unknown type")
	}

	params := stream.GetParameters()
	result := stream.AggrerateSoftwareManifest()

	format, saveToFile := params.OutputFormat, params.SaveToFile

	if len(saveToFile) > 0 {
		err := helper.SaveToFile(result, saveToFile, format.String())
		if err != nil {
			log.Fatal("Failed to save results to file :", err.Error())
		}
		return data
	}

	switch format {
	case types.Table:
		table.Show(table.Create(), duration)
	case types.JSON:
		json.DisplayResults(result)
	default:
		log.Fatal("Unknown output format")
	}

	return data
}
