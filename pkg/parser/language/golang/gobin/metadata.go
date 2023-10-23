package gobin

import (
	"runtime/debug"

	"github.com/carbonetes/diggity/pkg/model/metadata"
)

// Initialize Go metadata values from content
func parseMetadata(buildInfo *debug.BuildInfo, module *debug.Module) *metadata.GoBinMetadata {
	var metadata metadata.GoBinMetadata

	if buildInfo == nil {
		return nil
	}

	if module != nil {
		if len(buildInfo.Settings) > 0 {
			metadata.Architecture = parseBuildSettings(buildInfo.Settings, "GOARCH")
			metadata.Compiler = parseBuildSettings(buildInfo.Settings, "-compiler")
			metadata.OS = parseBuildSettings(buildInfo.Settings, "GOOS")
		}
		metadata.GoCompileRelease = buildInfo.GoVersion
		metadata.H1Digest = buildInfo.Main.Sum
		metadata.Path = buildInfo.Main.Path
		metadata.Version = buildInfo.Main.Version

		return &metadata
	}

	metadata.GoCompileRelease = buildInfo.GoVersion
	metadata.Path = buildInfo.Main.Path
	metadata.Version = buildInfo.Main.Version

	return &metadata
}
