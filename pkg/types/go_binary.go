package types

import "runtime/debug"

type GoBinary struct {
	File      string
	BuildInfo *debug.BuildInfo
}
