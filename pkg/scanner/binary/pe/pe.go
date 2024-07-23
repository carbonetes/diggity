package pe

import (
	"github.com/saferwall/pe"
)

type PEFile struct {
	Path string
	File *pe.File
}

func ParsePE(data []byte, path string) (*PEFile, bool) {

	file, err := pe.NewBytes(data, &pe.Options{})
	if err != nil {
		return nil, false
	}

	err = file.Parse()
	if err != nil {
		return nil, false
	}

	return &PEFile{
		path,
		file,
	}, true
}
