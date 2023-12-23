package helper

import "gopkg.in/yaml.v3"

func ToYAML(t interface{}) ([]byte, error) {
	return yaml.Marshal(t)
}

func WriteYAML(t interface{}, path string) error {
	b, err := ToYAML(t)
	if err != nil {
		return err
	}
	return WriteFile(b, path)
}

