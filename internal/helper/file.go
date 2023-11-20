package helper

import (
	"fmt"
	"os"
	"strings"
)

// IsDirExists checks if a directory exists and is valid.
func IsDirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// IsDir checks if the given input is a directory.
// It returns true if the input is a directory, false otherwise.
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

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

func SaveToFile(data interface{}, path, extension string) error {
	path = AddFileExtension(path, extension)
	// err := os.MkdirAll(filepath.Dir(path), 0700)
	// if err != nil {
	// 	return err
	// }
	// out, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	// if err != nil {
	// 	return err
	// }
	// defer out.Close()

	switch extension {
	case "json":
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
	default:
		return fmt.Errorf("invalid file extension: %s", extension)
	}

	return nil
}

func AddFileExtension(filename, outputType string) string {
	lastDotIndex := strings.LastIndex(filename, ".")
	if lastDotIndex != -1 {
		filename = filename[:lastDotIndex]
	}
	switch outputType {
	case "json", "cdx-json":
		return filename + ".json"
	case "xml", "cdx-xml":
		return filename + ".xml"
	default:
		return filename + ".txt"
	}
}
