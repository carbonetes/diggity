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

func Done() {
	p.Send(resultMsg{done: true})
}