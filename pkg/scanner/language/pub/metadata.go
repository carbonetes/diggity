package pub

import (
	"github.com/carbonetes/diggity/internal/log"
	"gopkg.in/yaml.v3"
)

type Lockfile struct {
	Packages map[string]Package `yaml:"packages"`
}
type Package struct {
	Dependency  string      `yaml:"dependency"`
	Description Description `yaml:"description"`
	Source      string      `yaml:"source"`
	Version     string      `yaml:"version"`
}
type Description struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

func readManifestFile(content []byte) map[string]interface{} {
	var metadata map[string]interface{}
	err := yaml.Unmarshal(content, &metadata)
	if err != nil {
		return nil
	}
	return metadata
}

func readLockFile(content []byte) *Lockfile {
	var metadata Lockfile
	err := yaml.Unmarshal(content, &metadata)
	if err != nil {
		log.Debug("Failed to unmarshal pubspec.lock")
		return nil
	}
	return &metadata
}
