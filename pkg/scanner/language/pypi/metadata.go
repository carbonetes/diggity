package pypi

import (
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/carbonetes/diggity/internal/log"
)

type PoetryLock struct {
	Packages []Package `json:"package"`
}

type Package struct {
	Name           string                `json:"name"`
	Version        string                `json:"version"`
	Description    string                `json:"description"`
	Optional       bool                  `json:"optional"`
	PythonVersions string                `json:"python-versions"`
	Files          []File                `json:"files"`
	Dependencies   map[string]Dependency `json:"dependencies"`
	Extras         map[string][]string   `json:"extras"`
}

type File struct {
	File string `json:"file"`
	Hash string `json:"hash"`
}

type Dependency struct {
	Version  string `json:"version,omitempty"`
	Optional bool   `json:"optional,omitempty"`
	Markers  string `json:"markers,omitempty"`
}

func readManifestFile(content []byte) map[string]interface{} {
	metadata := make(map[string]interface{})

	var key, value, prev string
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.Contains(line, ": ") {
			keyvalue := strings.SplitN(line, ": ", 2)
			if len(keyvalue) != 2 {
				continue
			}
			key, value = strings.TrimSpace(keyvalue[0]), strings.TrimSpace(keyvalue[1])
		} else {
			value = strings.TrimSpace(value + line)
			key = prev
		}

		if len(value) > 0 && value != " " {
			value = strings.Replace(value, "\r\n", "", -1)
			value = strings.Replace(value, "\r", "", -1)
			metadata[key] = strings.TrimSpace(value)
		}

	}
	return metadata
}

func readRequirementsFile(content []byte) [][]string {
	var attributes [][]string

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		props := strings.SplitAfterN(line, "==", 2)
		if len(props) != 2 {
			continue
		}
		attributes = append(attributes, props)
	}

	return attributes
}

func readPoetryLockFile(content []byte) PoetryLock {
	var lockFile PoetryLock
	if _, err := toml.Decode(string(content), &lockFile); err != nil {
		log.Error("Failed to decode poetry.lock file")
	}
	return lockFile
}