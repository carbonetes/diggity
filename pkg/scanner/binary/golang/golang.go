package golang

import (
	"debug/buildinfo"
	"io"
)

func Parse(r io.ReaderAt) (*buildinfo.BuildInfo, bool) {
	build, err := buildinfo.Read(r)
	if err != nil {
		return nil, false
	}

	return build, true
}
