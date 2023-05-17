package distro

import (
	"bufio"
	"os"
	"strings"
)

type Metadata map[string]string

func parseMetadata(path string) (*Metadata, error) {
	metadata := make(Metadata)
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var key, value string
		keyValue := strings.TrimSpace(scanner.Text())
		if strings.Contains(keyValue, "=") {
			keyValues := strings.SplitN(keyValue, "=", 2)
			key = keyValues[0]
			value = keyValues[1]
		}
		if len(key) > 0 && key != " " {
			value = strings.Replace(value, "\r\n", "", -1)
			value = strings.ReplaceAll(value, "\"", "")
			metadata[key] = strings.Replace(value, "\r ", "", -1)
			metadata[key] = strings.TrimSpace(metadata[key])
		}
	}
	return &metadata, nil
}
