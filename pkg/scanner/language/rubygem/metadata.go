package rubygem

import (
	"strings"
)

func readManifestFile(content []byte) [][]string {
	var attributes [][]string

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) < 5 {
			continue
		}

		if strings.Count(line[:5], " ") != 4 {
			continue
		}

		props := strings.Fields(line)

		if len(props) != 2 {
			continue
		}
		attributes = append(attributes, props)
	}

	return attributes
}

func readGemspecFile(content []byte) map[string]interface{} {
	metadata := make(map[string]interface{})
	lines := strings.Split(string(content), "\n")
	var key, value, prev string
	for _, line := range lines {
		if strings.Contains(line, "=") {
			keyvalue := strings.SplitN(line, "=", 2)
			if len(keyvalue) != 2 {
				continue
			}
			key, value = strings.TrimSpace(keyvalue[0]), strings.TrimSpace(keyvalue[1])
			if strings.Contains(value, "%") || strings.Contains(value, "if Gem") {
				value = ""
			}
		} else {
			value = strings.TrimSpace(value + line)
			key = prev
		}

		if len(value) > 0 && value != " " {
			value = strings.ReplaceAll(value, " ", "")
			value = strings.Replace(value, ".s", "", -1)
			key = strings.Replace(key, "\r\n", "", -1)
			key = strings.ReplaceAll(key, ".freeze", "")
			metadata[key] = strings.TrimSpace(value)
		}
		prev = value
	}

	return metadata
}
