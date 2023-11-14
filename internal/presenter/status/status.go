package status

import (
	"github.com/carbonetes/diggity/internal/logger"
)

var log = logger.GetLogger()

func ScanFile(data interface{}) interface{} {
	file, ok := data.(string)
	if !ok {
		log.Error("ScanFile received unknown type")
	}
	p.Send(resultMsg{file: file, done: false})
	return data
}

func ScanCompleteStatus(data interface{}) interface{} {
	_, ok := data.(bool)
	if !ok {
		log.Error("ScanComplete received unknown type")
	}
	p.Send(resultMsg{done: true})
	return data
}
