package conan

import (
	"encoding/json"
	"strings"

	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
)

// ConanMetadata conan metadata
type Metadata struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ConanLockMetadata conan.lock metadata
type LockMetadata struct {
	GraphLock GraphLock `json:"graph_lock"`
	Version   string    `json:"version"`
}

// GraphLock conan.lock metadata containing nodes
type GraphLock struct {
	Nodes            map[string]LockNode `json:"nodes"`
	RevisionsEnabled bool                `json:"revisions_enabled"`
}

// ConanLockNode conan.lock packages metadata
type LockNode struct {
	Ref            string      `json:"ref"`
	Path           string      `json:"path,omitempty"`
	Context        string      `json:"context,omitempty"`
	Requires       []string    `json:"requires,omitempty"`
	PackageID      string      `json:"package_id,omitempty"`
	Prev           string      `json:"prev,omitempty"`
	BuildRequires  string      `json:"build_requires,omitempty"`
	PythonRequires string      `json:"py_requires,omitempty"`
	Options        interface{} `json:"options,omitempty"`
}

func readManifestFile(content []byte) []Metadata {
	var packages []Metadata
	attributes := helper.SplitContentsByEmptyLine(string(content))
	for _, attribute := range attributes[1:] {
		if len(attribute) == 0 {
			continue
		}
		if !strings.Contains(attribute, "[requires]") {
			continue
		}

		lines := strings.Split(attribute, "\n")
		for _, line := range lines[1:] {
			if !strings.Contains(line, "/") {
				continue
			}
			if !strings.Contains(line, "[") {
				continue
			}

			props := strings.Split(line, "/")
			if len(props) < 1 {
				continue
			}
			name, version := props[0], props[1]
			packages = append(packages, Metadata{Name: name, Version: version})
		}

	}
	return packages
}

func readLockFile(content []byte) *LockMetadata {
	var metadata LockMetadata
	err := json.Unmarshal(content, &metadata)
	if err != nil {
		log.Debug("Failed to unmarshal conan.lock")
		return nil
	}

	return &metadata
}
