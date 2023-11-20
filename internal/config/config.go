package config

import (
	"github.com/carbonetes/diggity/internal/scanner"
	"github.com/carbonetes/diggity/pkg/types"
)

type Config struct {
	OutputFormat     string                   `json:"output_format" yaml:"output_format"`
	Quiet            bool                     `json:"quiet" yaml:"quiet"`
	Scanners         []string                 `json:"scanners" yaml:"scanners"`
	AllowFileListing bool                     `json:"allow_file_listing" yaml:"allow_file_listing"`
	AllowScanSecret  bool                     `json:"allow_scan_secret" yaml:"allow_scan_secret"`
	MaxFileSize      int64                    `json:"max_file_size" yaml:"max_file_size"`
	Registry         types.RegistryParameters `json:"registry" yaml:"registry"`
}

func New() Config {
	return Config{
		OutputFormat:     types.Table.String(),
		Quiet:            false,
		Scanners:         scanner.All,
		AllowFileListing: false,
		MaxFileSize:      10485760,
		Registry: types.RegistryParameters{
			URI:      "",
			Username: "",
			Password: "",
			Token:    "",
		},
	}
}
