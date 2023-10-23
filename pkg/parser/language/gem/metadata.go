package gem

import (
	"bufio"
	"os"
	"strings"
)

// Metadata  metadata
type Metadata map[string]interface{}

func parseMetadata(path string) (*Metadata, error) {
	gemFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer gemFile.Close()

	scanner := bufio.NewScanner(gemFile)

	var value string
	var attribute string
	var previousAttribute string

	metadata := make(Metadata)

	for scanner.Scan() {
		keyValue := scanner.Text()

		if strings.Contains(keyValue, "=") {
			keyValues := strings.SplitN(keyValue, "=", 2)
			attribute = keyValues[0]
			value = keyValues[1]

			//check if attribute is invalid - set to null if invalid
			if strings.Contains(attribute, "%") || strings.Contains(attribute, "if Gem") {
				//clear attribute
				attribute = ""
			}
		} else {
			value = strings.TrimSpace(value + keyValue)
			attribute = previousAttribute
		}

		if len(attribute) > 0 && attribute != " " {
			attribute = strings.ReplaceAll(attribute, " ", "")
			attribute = strings.Replace(attribute, "s.", "", -1)
			value = strings.Replace(value, "\r\n", "", -1)
			value = strings.ReplaceAll(value, ".freeze", "")
			metadata[attribute] = strings.ReplaceAll(value, "\"", "")
			metadata[attribute] = strings.TrimSpace(metadata[attribute].(string))
		}

		previousAttribute = attribute
	}
	return &metadata, nil
}
