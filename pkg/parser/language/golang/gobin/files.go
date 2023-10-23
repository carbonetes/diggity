package gobin

import (
	"debug/buildinfo"
	"os"
	"runtime/debug"
	"strings"

	"github.com/carbonetes/diggity/pkg/parser/language/golang"
)

func readFile(path string) (*debug.BuildInfo, error) {
	// Modify file permissions to allow read
	err := os.Chmod(path, 0777)
	if err != nil {
		if strings.Contains(err.Error(), golang.NoFileErrWin) || strings.Contains(err.Error(), golang.NoFileErrMac) {
			return nil, nil
		}
		return nil, err
	}

	binFile, err := os.Open(path)
	if err != nil {
		if strings.Contains(err.Error(), golang.NoFileErrWin) || strings.Contains(err.Error(), golang.NoFileErrMac) {
			return nil, nil
		}
		return nil, err
	}
	defer binFile.Close()

	buildData, err := buildinfo.Read(binFile)

	// Check if file is Go bin
	if err != nil {
		// Handle expected errors
		if err.Error() == "unrecognized file format" ||
			err.Error() == "not a Go executable" ||
			err.Error() == "EOF" {
			return nil, nil
		}
		return nil, err
	}
	return buildData, nil
}
