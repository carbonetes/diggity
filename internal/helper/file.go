package helper

import (
	"fmt"
	"os"
	"strings"
)



// IsFileExists checks if a file exists and is valid.
func IsFileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return !info.IsDir(), nil
}

func SaveToFile(data interface{}, path, format string) error {
	path = AddFileExtension(path, format)
	switch format {
	case "json", "cdx-json", "spdx-json":
		jsonData, err := ToJSON(data)
		if err != nil {
			return err
		}

		err = os.WriteFile(path, jsonData, 0644)
		if err != nil {
			return err
		}
	case "yaml":
		yamlData, err := ToYAML(data)
		if err != nil {
			return err
		}
		err = os.WriteFile(path, yamlData, 0644)
		if err != nil {
			return err
		}
	case "xml", "spdx-xml":
		xmlData, err := ToXML(data)
		if err != nil {
			return err
		}
		err = os.WriteFile(path, xmlData, 0644)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid format: %s", format)
	}

	fmt.Printf("File saved to %s", path)

	return nil
}

func AddFileExtension(filename, outputType string) string {
	lastDotIndex := strings.LastIndex(filename, ".")
	if lastDotIndex != -1 {
		filename = filename[:lastDotIndex]
	}
	switch outputType {
	case "json", "cdx-json", "spdx-json":
		return filename + ".json"
	case "xml", "cdx-xml":
		return filename + ".xml"
	default:
		return filename + ".txt"
	}
}

func WriteFile(data []byte, path string) error {
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
