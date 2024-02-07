package linux

import (
	"strings"

	"github.com/carbonetes/diggity/pkg/types"
)

// Parsing different OS release files
func parse(manifest types.ManifestFile) types.OSRelease {
	if manifest.Path == "etc/debian_version" {
		return parseDebianVersion(manifest)
	}

	if manifest.Path == "etc/centos-release" {
		return parseCentOSRelease(manifest)
	}

	if manifest.Path == "etc/redhat-release" {
		return parseRedhatRelease(manifest)
	}

	lines := strings.Split(string(manifest.Content), "\n")
	metadata := make(map[string]interface{})
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(parts[0])
		value := strings.Trim(parts[1], "\"") // Remove surrounding quotes
		metadata[key] = value
	}
	return types.OSRelease{
		File:    manifest.Path,
		Release: metadata,
	}
}

func parseDebianVersion(manifest types.ManifestFile) types.OSRelease {
	return types.OSRelease{
		File:    "etc/debian_version",
		Release: map[string]interface{}{"version": string(manifest.Content)},
	}
}

func parseCentOSRelease(manifest types.ManifestFile) types.OSRelease {
	return types.OSRelease{
		File:    "etc/centos-release",
		Release: map[string]interface{}{"name": string(manifest.Content)},
	}
}

func parseRedhatRelease(manifest types.ManifestFile) types.OSRelease {
	return types.OSRelease{
		File:    "etc/redhat-release",
		Release: map[string]interface{}{"name": string(manifest.Content)},
	}
}
