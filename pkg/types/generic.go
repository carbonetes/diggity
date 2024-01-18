package types

import (
	"bytes"
	"debug/elf"
)

type Generic struct {
	Path   string
	Name   string
	File   *elf.File
	ROData []string
}

func (g *Generic) ReadROData() {
	// Read the .rodata section, which usually contains strings.
	sec := g.File.Section(".rodata")
	if sec == nil {
		return
	}

	data, err := sec.Data()
	if err != nil {
		return
	}

	// Split the data into null-terminated strings.
	var strings []string
	for _, b := range bytes.Split(data, []byte{0}) {
		if len(b) > 0 {
			strings = append(strings, string(b))
		}
	}
	g.ROData = strings
}
