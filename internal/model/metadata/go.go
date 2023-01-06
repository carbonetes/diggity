package metadata

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
