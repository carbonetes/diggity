package docker

import (
	"os"
	"strings"
)

func getTarDir(dir, argDir string) (*os.File, error) {
	tarDirectory, err := os.Open(dir)
	if err != nil {
		if len(argDir) > 0 {
			tarDirectory, err = os.Open(argDir)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return tarDirectory, nil
}

// Get JSON files from extracted image
func getJSONFilesFromDir(root string) ([]string, error) {
	var files []string
	fileInfo, err := os.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		if strings.HasSuffix(file.Name(), ".json") {
			files = append(files, root+string(os.PathSeparator)+file.Name())
		}
	}
	return files, nil
}
