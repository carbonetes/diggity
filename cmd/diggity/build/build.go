package build

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const placeholder = "not available"

// Version - Build Information
type Version struct {
	// Version
	AppName   string `json:"appName"`
	Version   string `json:"version"`
	BuildDate string `json:"buildDate"`
	// Git
	GitCommit string `json:"gitCommit"`
	GitDesc   string `json:"gitDesc"`
	// Golang
	GoVersion string `json:"goVersion"`
	Compiler  string `json:"compiler"`
	Platform  string `json:"platform"`
}

// Build-time arguments, no values as default
var (
	appName, _  = os.Executable()
	baseAppName = filepath.Base(appName)
	version     = placeholder
	buildDate   = placeholder
	gitCommit   = placeholder
	gitDesc     = placeholder
	platform    = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
)

// FromBuild provides all version details
func FromBuild() Version {
	return Version{
		// Version
		AppName:   baseAppName,
		Version:   version,
		BuildDate: buildDate,
		// Git
		GitCommit: gitCommit,
		GitDesc:   gitDesc,
		// Golang
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  platform,
	}
}