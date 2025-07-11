package version

import (
	"fmt"
	"runtime"
)

// Build information set at compile time
var (
	// Version is the current version of diggity
	Version = "dev"

	// GitCommit is the git commit hash
	GitCommit = "unknown"

	// GitBranch is the git branch
	GitBranch = "unknown"

	// BuildDate is the date when the binary was built
	BuildDate = "unknown"

	// BuildUser is the user who built the binary
	BuildUser = "unknown"
)

// Info contains version and build information
type Info struct {
	Version     string `json:"version"`
	GitCommit   string `json:"git_commit"`
	GitBranch   string `json:"git_branch"`
	BuildDate   string `json:"build_date"`
	BuildUser   string `json:"build_user"`
	GoVersion   string `json:"go_version"`
	Compiler    string `json:"compiler"`
	Platform    string `json:"platform"`
	Application string `json:"application"`
}

// GetInfo returns version and build information
func GetInfo() Info {
	return Info{
		Version:     Version,
		GitCommit:   GitCommit,
		GitBranch:   GitBranch,
		BuildDate:   BuildDate,
		BuildUser:   BuildUser,
		GoVersion:   runtime.Version(),
		Compiler:    runtime.Compiler,
		Platform:    fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Application: "diggity",
	}
}

// String returns a formatted version string
func (i Info) String() string {
	return fmt.Sprintf("%s %s (%s)", i.Application, i.Version, i.GitCommit[:8])
}

// DetailedString returns a detailed version string with all build information
func (i Info) DetailedString() string {
	return fmt.Sprintf(`%s %s
Git Commit: %s
Git Branch: %s  
Build Date: %s
Build User: %s
Go Version: %s
Compiler:   %s
Platform:   %s`,
		i.Application,
		i.Version,
		i.GitCommit,
		i.GitBranch,
		i.BuildDate,
		i.BuildUser,
		i.GoVersion,
		i.Compiler,
		i.Platform,
	)
}

// GetVersion returns just the version string
func GetVersion() string {
	return Version
}

// GetVersionInfo returns a formatted version info string
func GetVersionInfo() string {
	info := GetInfo()
	return info.String()
}

// GetDetailedVersionInfo returns detailed version information
func GetDetailedVersionInfo() string {
	info := GetInfo()
	return info.DetailedString()
}
