package golang

import (
	"github.com/carbonetes/diggity/internal/log"
	"golang.org/x/mod/modfile"
)

// GoBinMetadata go binary metadata
type GoBinMetadata struct {
	Architecture     string `json:"architecture,omitempty"`
	Compiler         string `json:"compiler,omitempty"`
	OS               string `json:"os,omitempty"`
	GoCompileRelease string `json:"goCompileRelease,omitempty"`
	H1Digest         string `json:"h1Digest,omitempty"`
	Path             string `json:"path,omitempty"`
	Version          string `json:"version,omitempty"`
}

// GoModMetadata go module metadata
type GoModMetadata struct {
	Path    string `json:"path"`
	Version string `json:"version"`
}

// GoDevelMetadata - go devel metadata
type GoDevelMetadata struct {
	Path             string `json:"path"`
	Version          string `json:"version"`
	GoCompileRelease string `json:"goCompileRelease"`
}

func readManifestFile(content []byte, path string) *modfile.File {
	modFile, err := modfile.Parse(path, content, nil)
	if err != nil || modFile == nil {
		log.Debug("Failed to parse go.mod file")
		return nil
	}

	return modFile
}
