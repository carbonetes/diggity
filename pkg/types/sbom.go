package types

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/carbonetes/diggity/internal/helper"
	"github.com/docker/distribution/version"
)

const SchemaVersion = "1.0"

type Metadata struct {
	Author    string    `json:"author"`
	Timestamp time.Time `json:"timestamp"`
	Tool      string    `json:"tool"`
}

type SBOM struct {
	Schema     string `json:"schema"`
	Version    string `json:"version"`
	Serial     string `json:"serial"`
	Metadata   `json:"metadata"`
	Components []Component `json:"components"`
}

func NewSBOM() SBOM {
	return SBOM{
		Schema:  "https://github.com/carbonetes/diggity/schema/sbom/json/schema-1.0.json",
		Serial:  helper.GenerateURN(MakeNid()),
		Version: SchemaVersion,
		Metadata: Metadata{
			Author:    "diggity@" + version.Version,
			Timestamp: time.Now(),
			Tool:      "github.com/carbonetes/diggity",
		},
	}
}

func (c *Component) AddDependency(dependency Component) {
	c.Dependencies = append(c.Dependencies, Dependency{
		ParentID: c.ID,
		ChildID:  dependency.ID,
	})
}

func (c Component) ToJSON() string {
	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func (s SBOM) ToJSON() string {
	data, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func MakeNid() string {
	return fmt.Sprintf("diggity-schema-%s", SchemaVersion)
}
