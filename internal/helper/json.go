package helper

import (
	"bytes"
	"encoding/json"
)

func ToJSON(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", " ")
	err := encoder.Encode(t)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
