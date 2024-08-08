package linux

import (
	"strings"

	"github.com/carbonetes/diggity/pkg/types"
)

// Parsing different OS release files
func parse(file types.ManifestFile) types.OSRelease {
	if file.Path == "etc/debian_version" {
		return parseDebianVersion(file)
	}

	if file.Path == "etc/centos-release" {
		return parseCentOSRelease(file)
	}

	if file.Path == "etc/redhat-release" {
		return parseRedhatRelease(file)
	}

	lines := strings.Split(string(file.Content), "\n")
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
		File:    file.Path,
		Release: metadata,
	}
}

func parseDebianVersion(file types.ManifestFile) types.OSRelease {
	return types.OSRelease{
		File:    "etc/debian_version",
		Release: map[string]interface{}{"version": string(file.Content)},
	}
}

func parseCentOSRelease(file types.ManifestFile) types.OSRelease {
	return types.OSRelease{
		File:    "etc/centos-release",
		Release: map[string]interface{}{"name": string(file.Content)},
	}
}

func parseRedhatRelease(file types.ManifestFile) types.OSRelease {
	return types.OSRelease{
		File:    "etc/redhat-release",
		Release: map[string]interface{}{"name": string(file.Content)},
	}
}
