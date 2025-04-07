package gradle

import (
	"github.com/pelletier/go-toml"
)

type Metadata struct {
	Vendor  string `json:"vendor"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

// type Library map[string]Metadata
type VersionMetadata struct {
	Libraries map[string]Library  `toml:"libraries,omitempty" json:"libraries,omitempty"`
	Versions  map[string]string   `toml:"versions,omitempty" json:"versions,omitempty"`
	Bundles   map[string][]string `toml:"bundles,omitempty" json:"bundles,omitempty"`
	Plugins   map[string]Plugin   `toml:"plugins,omitempty" json:"plugins,omitempty"`
}

type Library struct {
	Module  *string `toml:"module,omitempty" json:"module,omitempty"`
	Name    *string `toml:"name,omitempty" json:"name,omitempty"`
	Group   *string `toml:"group,omitempty" json:"group,omitempty"`
	Version any     `toml:"version,omitempty" json:"version,omitempty"`
}

type Plugin struct {
	ID      string `toml:"id,omitempty" json:"id,omitempty"`
	Version any    `toml:"version,omitempty" json:"version,omitempty"`
}

func read(content string) *VersionMetadata {
	var data VersionMetadata
	err := toml.Unmarshal([]byte(content), &data)
	if err != nil {
		panic(err)
	}

	return &data
}
