package status

import (
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/stream"
)

func FileListWatcher(data interface{}) interface{} {
	file, ok := data.(string)
	if !ok {
		log.Error("FileCheckWatcher received unknown type")
	}
	p.Send(resultMsg{file: file, done: false})
	stream.AddFile(file)
	return data
}

func ScanElapsedStoreWatcher(data interface{}) interface{} {
	_, ok := data.(float64)
	if !ok {
		log.Error("ScanElapsedStoreWatcher received unknown type")
	}
	p.Send(resultMsg{done: true})
	return data
}

func ErrorStoreWatcher(data interface{}) interface{} {
	err, ok := data.(error)
	if !ok {
		log.Error("ErrorStoreWatcher received unknown type")
	}
	p.Send(errorMsg{err: err, quit: false})
	return data
}
