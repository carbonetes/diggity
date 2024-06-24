package npm

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/carbonetes/diggity/internal/log"
	"gopkg.in/yaml.v3"
)

// PackageJSON - packages.json model
type Metadata struct {
	Name            string                 `json:"name"`
	Version         string                 `json:"version"`
	Latest          []string               `json:"latest"`
	Contributors    interface{}            `json:"contributors"`
	License         interface{}            `json:"license"`
	Homepage        string                 `json:"homepage"`
	Description     string                 `json:"description"`
	Dependencies    map[string]interface{} `json:"dependencies"`
	DevDependencies map[string]interface{} `json:"devDependencies"`
	Repository      interface{}            `json:"repository"`
	Author          interface{}            `json:"author"`
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

func readManifestFile(content []byte) *Metadata {
	var metadata Metadata
	err := json.Unmarshal(content, &metadata)
	if err != nil {
		log.Debug("Failed to unmarshal package.json")
		return nil
	}
	return &metadata
}

func readPackageLockfile(content []byte) *PackageLock {
	var metadata PackageLock
	err := json.Unmarshal(content, &metadata)
	if err != nil {
		log.Debug("Failed to unmarshal package-lock.json")
		return nil
	}

	return &metadata
}

func readPnpmLockfile(content []byte) *PnpmLockfile {
	var metadata PnpmLockfile
	err := yaml.Unmarshal(content, &metadata)
	if err != nil {
		log.Debug("Failed to unmarshal pnpm-lock.yaml")
		return nil
	}
	return &metadata
}

type Package struct {
	Version      string
	Resolution   string
	Dependencies map[string]string
}

type YarnLockfile struct {
	Version      string            `yaml:"version"`
	Resolution   string            `yaml:"resolution"`
	Dependencies map[string]string `yaml:"dependencies"`
	Checksum     string            `yaml:"checksum"`
	LanguageName string            `yaml:"languageName"`
	LinkType     string            `yaml:"linkType"`
}

func parseYarnLock(content []byte) (map[string]YarnLockfile, error) {
	packages := make(map[string]YarnLockfile)
	r := bytes.NewReader(content)
	data, err := io.ReadAll(r)
	if err != nil {
		log.Debug("Failed to read yarn.lock")
		return nil, err
	}

	err = yaml.Unmarshal(data, &packages)
	if err != nil {
		log.Debug("Failed to unmarshal yarn.lock")
		return nil, err
	}

	return packages, nil
}
