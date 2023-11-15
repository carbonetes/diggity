package presenter

import (
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
		log.Error("DisplayResults received unknown type")
	}

	format := stream.GetParameters().OutputFormat

	switch format {
	case types.Table:
		table.Show(table.Create(), duration)
	case types.JSON:
		json.DisplayResults(stream.AggrerateSoftwareManifest())
	default:
		log.Error("Unknown output format")
	}

	return data
}
