package hex

import (
	"regexp"
	"strings"
)

type HexMetadata struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	PkgHash    string `json:"pkgHash,omitempty"`
	PkgHashExt string `json:"pkgHashExt,omitempty"`
}

var (
	rebarLockRegex = regexp.MustCompile(`[\[{<">},: \]\n]+`)
	mixLockRegex   = regexp.MustCompile(`[%{}\n" ,:]+`)
)

func readRebarFile(content []byte) []HexMetadata {
	var packages []HexMetadata
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		token := rebarLockRegex.Split(line, -1)
		if len(token) != 7 {
			continue
		}
		name, version := token[1], token[4]
		metadata := HexMetadata{
			Name:    name,
			Version: version,
		}
		packages = append(packages, metadata)
	}
	return packages
}

func readMixFile(content []byte) []HexMetadata {
	var packages []HexMetadata
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		tokens := mixLockRegex.Split(line, -1)
		if len(tokens) != 6 {
			continue
		}
		name, version, hash, hashExt := tokens[1], tokens[4], tokens[5], tokens[len(tokens)-2]
		metadata := HexMetadata{
			Name:       name,
			Version:    version,
			PkgHash:    hash,
			PkgHashExt: hashExt,
		}
		packages = append(packages, metadata)
	}
	return packages
}
