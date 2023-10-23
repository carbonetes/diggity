package python

import (
	"regexp"
	"strings"
)

// Parse poetry file
func poetryFileMetadata(file string) map[string]string {
	fileHash := make(map[string]string)
	r := regexp.MustCompile(`"(.*?)"`)

	for _, fh := range r.FindAllString(file, -1) {
		// assign to hash if contains sha256
		if strings.Contains(fh, fileHashKey) {
			fileHash["hash"] = strings.Replace(fh, `"`, "", -1)
		} else {
			fileHash["file"] = strings.Replace(fh, `"`, "", -1)
		}
	}

	return fileHash
}
