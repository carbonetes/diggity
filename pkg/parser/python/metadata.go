package python

import (
	"fmt"
	"os"
	"strings"
)

// Metadata  metadata
type Metadata map[string]interface{}

// Parse python metadata
func parseMetadataFiles(m Metadata, path string) error {
	var mapValue = map[string]interface{}{}
	var files []map[string]interface{}
	var finalValue = map[string]interface{}{}
	var pathValue, algorithm, algorithmValue, valueSize string
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("file does not exists: %s", path)
	}

	fileinfo, _ := os.ReadFile(path)

	lines := strings.Split(string(fileinfo), "\n")
	for _, line := range lines {

		if strings.Contains(line, ",") {
			continue
		}
		keyValues := strings.Split(line, ",")

		if len(keyValues) != 2 {
			continue
		}

		pathValue = keyValues[0]
		tmpValue := keyValues[1]
		if tmpValue != "" {
			tmpSplitValue := strings.Split(tmpValue, "=")
			algorithm = tmpSplitValue[0]
			algorithmValue = tmpSplitValue[1]
		}
		valueSize = keyValues[2]

		if pathValue != "" {
			finalValue["path"] = pathValue
			mapValue["value"] = algorithmValue
			mapValue["algorithm"] = algorithm

			//ignore digest if algorithm is blank
			if algorithm != "" {
				finalValue["digest"] = mapValue
			}
			//ignore valueSize if blank
			if valueSize != "" {
				valueSize = strings.Replace(valueSize, "\r", "", -1)
				finalValue["size"] = valueSize
			}
			files = append(files, finalValue)
			algorithm = ""
			algorithmValue = ""
		}

	}
	m["Files"] = files
	return nil
}

// Parse requirements metadata
func parseRequirements(req string) (name string, version string) {
	reqMetadata := strings.Split(req, "==")
	versionMetadata := strings.TrimSpace(reqMetadata[1])

	name = strings.TrimSpace(reqMetadata[0])
	version = strings.Split(versionMetadata, " ")[0]

	return name, version
}
