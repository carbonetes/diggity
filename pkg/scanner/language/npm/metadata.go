package npm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"strings"

	"github.com/carbonetes/diggity/internal/log"
	"gopkg.in/yaml.v3"
)

// PackageJSON - packages.json model
type Metadata struct {
	Version      string                 `json:"version"`
	Latest       []string               `json:"latest"`
	Contributors interface{}            `json:"contributors"`
	License      interface{}            `json:"license"`
	Name         string                 `json:"name"`
	Homepage     string                 `json:"homepage"`
	Description  string                 `json:"description"`
	Dependencies map[string]interface{} `json:"dependencies"`
	Repository   interface{}            `json:"repository"`
	Author       interface{}            `json:"author"`
}

// Contributors - PackageJSON contributors
type Contributors struct {
	Name     string `json:"name" mapstruct:"name"`
	Username string `json:"email" mapstruct:"username"`
	URL      string `json:"url" mapstruct:"url"`
}

// Repository - PackageJSON repository
type Repository struct {
	Type string `json:"type" mapstructure:"type"`
	URL  string `json:"url" mapstructure:"url"`
}

// PackageLock - PackageLock model
type PackageLock struct {
	Requires        bool `json:"requires"`
	LockfileVersion int  `json:"lockfileVersion"`
	Dependencies    map[string]LockDependency
}

// LockDependency - PackageLock dependencies
type LockDependency struct {
	Version   string `json:"version"`
	Resolved  string `json:"resolved"`
	Integrity string `json:"integrity"`
	Requires  map[string]string
}

type Dependency struct {
	Version      string            `yaml:"version"`
	Resolution   string            `yaml:"resolution"`
	Dependencies map[string]string `yaml:"dependencies"`
	Checksum     string            `yaml:"checksum"`
	LanguageName string            `yaml:"languageName"`
	LinkType     string            `yaml:"linkType"`
}

type YarnPackage struct {
    Version      string
    Resolution   string
    Dependencies map[string]string
}

type PnpmLockfile struct {
	LockFileVersion string                 `yaml:"lockfileVersion"`
	Dependencies    map[string]interface{} `yaml:"dependencies"`
	Packages        map[string]interface{} `yaml:"packages"`
}

func readManifestFile(content []byte) Metadata {
	var metadata Metadata
	err := json.Unmarshal(content, &metadata)
	if err != nil {
		log.Error("Failed to unmarshal package.json")
	}
	return metadata
}

func readPackageLockfile(content []byte) PackageLock {
	var metadata PackageLock
	err := json.Unmarshal(content, &metadata)
	if err != nil {
		log.Error("Failed to unmarshal package-lock.json")
	}

	return metadata
}

func readPnpmLockfile(content []byte) PnpmLockfile {
	var metadata PnpmLockfile
	err := yaml.Unmarshal(content, &metadata)
	if err != nil {
		log.Error("Failed to unmarshal pnpm-lock.yaml")
	}
	return metadata
}

type Package struct {
	Version      string
	Resolution   string
	Dependencies map[string]string
}

func ParseYarnLock(content []byte) (map[string]Package, error) {
	r := bytes.NewReader(content)
	scanner := bufio.NewScanner(r)
	packages := make(map[string]Package)

	var currentPackage string
	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		// If the line ends with a colon, it's a package name
		if strings.HasSuffix(line, ":") {
			currentPackage = strings.TrimSuffix(line, ":")
			packages[currentPackage] = Package{
				Dependencies: make(map[string]string),
			}
			continue
		}

		// If we're inside a package definition
		if currentPackage != "" {
			parts := strings.SplitN(line, " ", 2)
			key := strings.TrimSpace(parts[0])
			value := ""
			if len(parts) > 1 {
				value = strings.TrimSpace(parts[1])
			}

			switch key {
			case "version:":
				packages[currentPackage] = Package{
					Version:      strings.Trim(value, "\""),
					Dependencies: packages[currentPackage].Dependencies,
				}
			case "resolution:":
				packages[currentPackage] = Package{
					Version:      packages[currentPackage].Version,
					Resolution:   strings.Trim(value, "\""),
					Dependencies: packages[currentPackage].Dependencies,
				}
			case "dependencies:":
				// Handle dependencies in the next lines
				for scanner.Scan() {
					line := scanner.Text()
					if strings.TrimSpace(line) == "" {
						break
					}
					parts := strings.SplitN(line, " ", 2)
					depName := strings.TrimSpace(parts[0])
					depVersion := ""
					if len(parts) > 1 {
						depVersion = strings.TrimSpace(parts[1])
					}
					packages[currentPackage].Dependencies[depName] = depVersion
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return packages, nil
}
