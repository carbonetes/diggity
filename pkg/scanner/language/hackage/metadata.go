package hackage

import (
	"strings"

	"github.com/carbonetes/diggity/internal/log"
	"gopkg.in/yaml.v3"
)

// StackConfig stack.yaml metadata containing extra-deps
type StackConfig struct {
	ExtraDeps []string `yaml:"extra-deps"`
}

// StackLockConfig stack.yaml.lock metadata containing packages
type StackLockConfig struct {
	Packages  []StackPackages `yaml:"packages"`
	Snapshots []interface{}   `yaml:"snapshots"`
}

// StackPackages stack.yaml.lock packages metadata
type StackPackages struct {
	Original Hackage `yaml:"original"`
}

// Hackage metadata
type Hackage struct {
	Hackage string `yaml:"hackage"`
}

// HackageMetadata haskell packages metadata
type HackageMetadata struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	PkgHash     string `json:"pkgHash,omitempty"`
	Size        string `json:"size,omitempty"`
	Revision    string `json:"revision,omitempty"`
	SnapshotURL string `json:"snapshotURL,omitempty"`
}

func readStackConfigFile(content []byte) StackConfig {
	var metadata StackConfig
	err := yaml.Unmarshal(content, &metadata)
	if err != nil {
		log.Error("Failed to unmarshal stack.yaml")
	}
	return metadata
}

func readStackLockConfigFile(content []byte) StackLockConfig {
	var metadata StackLockConfig
	err := yaml.Unmarshal(content, &metadata)
	if err != nil {
		log.Error("Failed to unmarshal stack.yaml.lock")
	}
	return metadata
}

func readManifestFile(content []byte) []string {
	var packages []string
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.Contains(line, "any.") {
			var pkg string
			// Remove constraints field
			if strings.Contains(line, "constraints:") {
				pkg = strings.Replace(line, "constraints:", "", -1)
			} else {
				pkg = line
			}
			if nv := formatCabalPackage(pkg); nv != "" {
				packages = append(packages, formatCabalPackage(pkg))
			}
		}
	}
	return packages
}

// Parse Name, Version, PkgHash, Size, and Revision from extra-deps
func parseExtraDep(dep string) (name string, version string, pkgHash string, size string, rev string) {
	pkg := strings.Split(dep, "@")
	nv := strings.Split(pkg[0], "-")
	name = strings.Join(nv[0:len(nv)-1], "-")
	version = nv[len(nv)-1]

	if len(pkg) > 1 {
		// Parse pkgHash if sha256 is detected
		if strings.Contains(pkg[1], "sha256") {
			hs := strings.Split(pkg[1], ",")
			pkgHash = hs[0]
			size = hs[1]
		}
		// Parse revision if rev is detected
		if strings.Contains(pkg[1], "rev") {
			rev = pkg[1]
		}
	}

	return name, version, pkgHash, size, rev
}
