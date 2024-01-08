package spdx

import (
	"os"
	"strings"

	"github.com/carbonetes/diggity/pkg/converter/spdx/licenses"
	"github.com/google/uuid"
)

func CheckLicense(id string) string {
	licenseList := licenses.List[strings.ToLower(id)]
	return licenseList
}

func FormatPath(path string) string {
	pathSlice := strings.Split(path, string(os.PathSeparator))
	return strings.Join(pathSlice, "/")
}

func FormatAuthor(authorString string) string {
	author := []string{}

	// Check for empty author
	if strings.TrimSpace(authorString) == "" {
		return ""
	}

	authorDetails := strings.Split(authorString, " ")
	if len(authorDetails) == 1 {
		return authorDetails[0]
	}

	for _, detail := range authorDetails {
		if strings.Contains(detail, "http") && strings.Contains(detail, ".") && strings.Contains(detail, "/") {
			continue
		}
		author = append(author, detail)
	}

	return strings.Join(author, " ")
}

func FormatNamespace(input string) string {
	return namespace + input + "-" + uuid.NewString()
}
