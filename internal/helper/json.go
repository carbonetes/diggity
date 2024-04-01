package helper

import (
	"bytes"
	"encoding/json"
	"strings"
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

// CleanValue removes trailing carriage return, newline, and backslash from a JSON value.
func CleanValue(v interface{}) interface{} {
	switch v := v.(type) {
	case string:
		// Remove trailing carriage return, newline, and backslash
		return strings.NewReplacer("\r", "", "\n", "", "\\", "").Replace(v)
	case []interface{}:
		for i, u := range v {
			v[i] = CleanValue(u)
		}
	case map[string]interface{}:
		for k, u := range v {
			v[k] = CleanValue(u)
		}
	}
	return v
}

// CleanJSON removes trailing carriage return, newline, and backslash from a JSON string.
func CleanJSON(jsonString string) (string, error) {
	var v interface{}
	if err := json.Unmarshal([]byte(jsonString), &v); err != nil {
		return "", err
	}

	v = CleanValue(v)

	cleanedJSON, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(cleanedJSON), nil
}
