package nuget

import (
	"encoding/json"

	"github.com/carbonetes/diggity/internal/log"
)

// DotnetDeps - .NET Dependencies
type DotnetDeps struct {
	Libraries map[string]DotnetLibrary `json:"libraries"`
}

// DotnetLibrary - .NET libraries
type DotnetLibrary struct {
	Type     string `json:"type"`
	Path     string `json:"path"`
	Sha512   string `json:"sha512"`
	HashPath string `json:"hashPath"`
}

func readManifestFile(content []byte) DotnetDeps {
	var metadata DotnetDeps
	err := json.Unmarshal(content, &metadata)
	if err != nil {
		log.Error("Failed to unmarshal .deps.json")
	}
	return metadata
}
