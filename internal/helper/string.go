package helper

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var caser = cases.Title(language.English)

func ToTitle(str string) string {
	return caser.String(str)
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

func SplitContentsByEmptyLine(contents string) []string {
	attributes := regexp.
		MustCompile("\r\n").
		ReplaceAllString(contents, "\n")

	return regexp.
		MustCompile(`\n\s*\n`).
		Split(attributes, -1)
}

func GenerateURN(nid string) string {
	// Generate a new UUID
	uuid := uuid.New()

	// Construct the URN with the provided namespace identifier
	return fmt.Sprintf("urn:%s:%s", nid, uuid.String())
}

func SplitString(str string) []string {
	return regexp.
		MustCompile(`\s*[\s,;]+\s*`).
		Split(str, -1)
}

func SplitAndAppendStrings(target []string) []string {
	var result []string
	for _, str := range target {
		result = append(result, SplitString(str)...)
	}
	return result
}

func SplitAny(s string, seps string) []string {
	result := strings.FieldsFunc(s, func(r rune) bool {
		return strings.ContainsRune(seps, r)
	})
	if len(result) == 0 {
		return []string{s}
	}
	return result
}
