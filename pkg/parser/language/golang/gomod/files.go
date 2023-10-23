package gomod

import (
	"io"
	"os"
	"strings"

	"github.com/carbonetes/diggity/pkg/parser/language/golang"
	"golang.org/x/mod/modfile"
)

func readModFile(path string) (*modfile.File, error) {
	reader, err := os.Open(path)
	if err != nil {
		if strings.Contains(err.Error(), golang.NoFileErrWin) || strings.Contains(err.Error(), golang.NoFileErrMac) {
			return nil, nil
		}
		return nil, err
	}
	defer reader.Close()

	modContents, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	modFile, err := modfile.Parse(path, modContents, nil)
	if err != nil {
		return nil, err
	}

	return modFile, nil
}
