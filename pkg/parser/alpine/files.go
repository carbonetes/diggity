package alpine

import (
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
)

// Parse alpine files
func getAlpineFiles(content string) []model.File {

	var files []model.File
	keyValues := strings.Split(content, "\n")
	for idx := range keyValues {
		file := model.File{}

		// F: = File or Directory
		if strings.HasPrefix(keyValues[idx], "F:") {
			file.Path = strings.SplitN(keyValues[idx], ":", 2)[1]
			files = append(files, file)
		} else if strings.HasPrefix(keyValues[idx], "R:") {
			// reloop until F: or R: prefix is found
			file.Path = strings.SplitN(keyValues[idx], ":", 2)[1]
			for fileIdx := idx + 1; fileIdx < len(keyValues) && !strings.HasPrefix(keyValues[fileIdx], "R:"); {
				//  a:, M: = File Permissions
				if strings.HasPrefix(keyValues[fileIdx], "a:") || strings.HasPrefix(keyValues[fileIdx], "M:") {
					file.OwnerGID = strings.Split(keyValues[fileIdx], ":")[1]
					file.OwnerUID = strings.Split(keyValues[fileIdx], ":")[2]
					file.Permissions = strings.Split(keyValues[fileIdx], ":")[3]
				} else if /* Z: = Pull Checksum */ strings.HasPrefix(keyValues[fileIdx], "Z:") {
					digest := map[string]string{}
					digest["algorithm"] = "sha1"
					digest["value"] = strings.SplitN(keyValues[fileIdx], ":", 2)[1]
					file.Digest = digest
				}
				fileIdx++
			}
			files = append(files, file)
		}

	}

	return files
}
