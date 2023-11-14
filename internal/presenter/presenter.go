package presenter

import (
	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/internal/presenter/table"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

var log = logger.GetLogger()

func DisplayResults(data interface{}) interface{} {
	_, ok := data.(bool)
	if !ok {
		log.Error("ScanComplete received unknown type")
	}

	format := stream.GetParameters().OutputFormat

	switch format {
	case types.Table:
		table.Show(table.Create())
	default:
		log.Error("Unknown output format")
	}

	return data
}
