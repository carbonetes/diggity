package types

import "runtime/debug"

type GoBinary struct {
	File      string
	Path      string
	BuildInfo *debug.BuildInfo
}
