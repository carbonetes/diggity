package helper

import "gopkg.in/yaml.v3"

func ToYAML(t interface{}) ([]byte, error) {
	return yaml.Marshal(t)
}
