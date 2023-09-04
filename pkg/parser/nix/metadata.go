package nix

import (
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model/metadata"
)

var versionPattern = regexp.MustCompile(`-(?P<version>\d[a-zA-Z0-9]*(?:\.\d[a-zA-Z0-9]*){0,3}(?:-(?P<prerelease>\d*[.a-zA-Z-][.0-9a-zA-Z-]*)*)?(?:\+(?P<metadata>[.0-9a-zA-Z-]+(?:\.[.0-9a-zA-Z-]+)*))?)`)

func parseNixPath(input string) *metadata.NixMetadata {
	versionStart, version := findVersionWithPattern(input, versionPattern)
	if versionStart < 0 {
		return nil
	}
	hashName := strings.TrimSuffix(input[0:versionStart], "-")
	fields := strings.Split(hashName, "-")
	if len(fields) < 2 {
		return nil
	}
	hash, name := fields[0], strings.Join(fields[1:], "-")
	return &metadata.NixMetadata{
		Hash:    hash,
		Name:    name,
		Version: version,
	}
}

func findVersionWithPattern(input string, pattern *regexp.Regexp) (int, string) {
	match := pattern.FindAllStringSubmatchIndex(input, -1)
	if len(match) == 0 || len(match[0]) == 0 {
		return -1, ""
	}
	// TODO: check prerelease prefix and suffix match with regexp prerelease
	versionGroup := pattern.SubexpIndex("version")
	versionStart, versionStop := match[0][versionGroup*2], match[0][(versionGroup*2)+1]

	var version string
	if versionStart != -1 && versionStop != -1 {
		version = input[versionStart:versionStop]
	}

	return versionStart, version
}
