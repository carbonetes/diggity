package types

import (
	"encoding/json"
	"fmt"

	"time"

	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/logger"
	"github.com/docker/distribution/version"
)

const SchemaVersion = "1.0"

var log = logger.GetLogger()

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
	Total      int         `json:"total"`
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

func (c Component) ToJSON() string {
	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		log.Error(err)
	}
	return string(data)
}

func (s SBOM) ToJSON() string {
	data, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		log.Error(err)
	}
	return string(data)
}

func MakeNid() string {
	return fmt.Sprintf("diggity-schema-%s", SchemaVersion)
}
