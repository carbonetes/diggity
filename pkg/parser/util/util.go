package util

import (
	"os"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/pkg/model"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ParserNames slice of supported parser names
var ParserNames = []string{
	"apk",
	"deb",
	"java",
	"npm",
	"pnpm",
	"php",
	"python",
	"gem",
	"rpm",
	"dart",
	"nuget",
	"go",
	"rust-crate",
	"conan",
	"hackage",
	"pod",
	"hex",
	"portage",
	"alpm",
	"gradle",
}

var caser = cases.Title(language.English)

var log = logger.GetLogger()

// TrimUntilLayer Returns file path without layer hash
func TrimUntilLayer(location model.Location) string {
	directories := strings.Split(location.Path, string(os.PathSeparator))
	index := IndexOf(directories, location.LayerHash) + 1
	directory := ""
	for index < len(directories) {
		if index == len(directories)-1 {
			directory += directories[index]
		} else {
			directory += directories[index] + string(os.PathSeparator)
		}

		index++
	}
	directory = strings.ReplaceAll(directory, "\\", "/")
	return directory
}

// ParserEnabled checks if all or a specific parser is enabled
func ParserEnabled(parser string, enabledParsers *[]string) bool {
	if len(*enabledParsers) == 0 {
		return true
	}
	if StringSliceContains(*enabledParsers, parser) {
		return true
	}
	return false
}

// IndexOf returns index of a string from a slice
func IndexOf(array []string, s string) int {
	for idx, a := range array {
		if s == a {
			return idx
		}
	}
	return -1
}

// StringSliceContains checks if a string slice contains specified string
func StringSliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// FormatLockKeyVal formats .lock Key Value Data String
func FormatLockKeyVal(kv string) string {
	trimmed := strings.TrimSpace(kv)
	return strings.Replace(trimmed, `"`, "", -1)
}

// CleanUp clears temp files
func CleanUp(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Error(err)
	}
}

func SplitContentsByEmptyLine(contents string) []string {
	attributes := regexp.
		MustCompile("\r\n").
		ReplaceAllString(contents, "\n")

	return regexp.
		MustCompile(`\n\s*\n`).
		Split(attributes, -1)
}

func ToTitle(str string) string {
	return caser.String(str)
}
