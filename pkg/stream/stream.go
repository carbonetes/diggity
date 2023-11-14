package stream

import (
	"gitlab.com/jhumel/grove"
)

var (
	hub   *grove.Grove
	store *grove.Store
)

func init() {
	hub = grove.New()
	store = grove.NewStore(hub)
	SetDefaultValues()
}

func Emit(event string, data interface{}) {
	hub.Emit(event, data)
}

func Attach(event string, handler grove.Handler) {
	hub.Attach(event, handler)
}

func GetHub() *grove.Grove {
	return hub
}

func GetStore() *grove.Store {
	return store
}
