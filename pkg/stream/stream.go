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

func Watch(key string, handler grove.Handler) {
	store.Watch(key, handler)
}

func Set(key string, value interface{}) {
	store.Set(key, value)
}

func Get(key string) (interface{}, bool) {
	return store.Get(key)
}
