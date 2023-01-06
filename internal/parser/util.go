package parser

import (
	"os"
	"strings"

	"github.com/carbonetes/diggity/internal/model"
)

// ParserNames slice of supported parser names
var ParserNames = []string{
	"apk",
	"debian",
	"java",
	"npm",
	"composer",
	"python",
	"gem",
	"rpm",
	"dart",
	"nuget",
	"go",
}

// TrimUntilLayer Returns file path without layer hash
func TrimUntilLayer(location model.Location) string {
	directories := strings.Split(location.Path, string(os.PathSeparator))
	index := indexOf(directories, location.LayerHash) + 1
	directory := ""
	for index < len(directories) {
		if index == len(directories)-1 {
			directory += directories[index]
		} else {
			directory += directories[index] + string(os.PathSeparator)
		}

		index++
	}
	return directory
}

// Checks if all or a specific parser is enabled
func parserEnabled(parser string) bool {
	if len(*Arguments.EnabledParsers) == 0 {
		return true
	}
	if stringSliceContains(*Arguments.EnabledParsers, parser) {
		return true
	}
	return false
}

func indexOf(array []string, s string) int {
	for idx, a := range array {
		if s == a {
			return idx
		}
	}
	return -1
}

func stringSliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func sourceIsDir() bool {
	return len(*Arguments.Dir) > 0
}
