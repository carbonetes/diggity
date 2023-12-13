package npm

import (
	"encoding/json"

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

type YarnLockfile struct {
	Metadata     map[string]int `yaml:"__metadata"`
	Dependencies map[string]Dependency
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

func readYarnLockfile(content []byte) YarnLockfile {
	var lockfile YarnLockfile
	err := yaml.Unmarshal(content, &lockfile)
	if err != nil {
		log.Errorf("Failed to unmarshal yarn.lock file: %s", err.Error())
	}

	return lockfile
}

func readPnpmLockfile(content []byte) PnpmLockfile {
	var metadata PnpmLockfile
	err := yaml.Unmarshal(content, &metadata)
	if err != nil {
		log.Error("Failed to unmarshal pnpm-lock.yaml")
	}
	return metadata
}
