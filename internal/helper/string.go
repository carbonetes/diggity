package helper

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
)

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
