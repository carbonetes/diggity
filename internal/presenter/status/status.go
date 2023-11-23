package status

import (
	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/pkg/stream"
)

var log = logger.GetLogger()

func FileListWatcher(data interface{}) interface{} {
	file, ok := data.(string)
	if !ok {
		log.Fatal("FileCheckWatcher received unknown type")
	}
	p.Send(resultMsg{file: file, done: false})
	stream.AddFile(file)
	return data
}

func ScanElapsedStoreWatcher(data interface{}) interface{} {
	_, ok := data.(float64)
	if !ok {
		log.Fatal("ScanElapsedStoreWatcher received unknown type")
	}
	p.Send(resultMsg{done: true})
	return data
}
