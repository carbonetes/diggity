package types

import (
	"fmt"
	"os"
	"strings"

	"github.com/carbonetes/diggity/internal/helper"
)

type ScanType int
type OutputFormat string

const (
	Image      ScanType = 1
	Tarball    ScanType = 2
	Filesystem ScanType = 3

	JSON          OutputFormat = "json"
	Table         OutputFormat = "table"
	CycloneDXJSON OutputFormat = "cdx-json"
	CycloneDXXML  OutputFormat = "cdx-xml"
	SPDXJSON      OutputFormat = "spdx-json"
	SPDXXML       OutputFormat = "spdx-xml"
	SPDXTag       OutputFormat = "spdx-tag"
	SnapshotJSON  OutputFormat = "snapshot-json"
)

// DefaultSecretExtensions contains a list of common file extensions containing secrets.
// Additional Reference: https://blog.gitguardian.com/top-10-file-extensions/
var DefaultSecretExtensions = []string{"env", "h", "so", "sec", "pem", "properties", "xml", "yml", "yaml", "json", "py", "js", "ts", "PHP"}

type Parameters struct {
	ScanType         ScanType
	Input            string
	OutputFormat     OutputFormat
	Quiet            bool
	Scanners         []string
	AllowFileListing bool
	AllowPullTimeout bool
	Secrets          SecretParameters
	Registry         RegistryParameters
	Provenance       string
}

type SecretParameters struct {
	AllowSearch       bool
	MaxFileSize       int64
	ContentRegex      string
	Extensions        []string
	ExcludedFilenames []string
}

type RegistryParameters struct {
	URI      string
	Username string
	Password string
	Token    string
}

func (o OutputFormat) String() string {
	return string(o)
}

func (p *Parameters) GetScanType() error {
	if strings.Contains(p.Input, ":") {
		p.ScanType = Image
	} else if strings.Contains(p.Input, ".tar") {
		p.ScanType = Tarball
	} else if strings.Contains(p.Input, string(os.PathSeparator)) {
		exists, err := helper.IsDirExists(p.Input)
		if err != nil {
			return err
		}
		if exists {
			p.ScanType = Filesystem
		}
	} else {
		return fmt.Errorf("Invalid input value %v", p.Input)
	}
	return nil
}

func DefaultParameters() Parameters {
	return Parameters{
		ScanType:         0,
		Input:            "",
		OutputFormat:     JSON,
		Quiet:            false,
		Scanners:         nil,
		AllowFileListing: false,
		AllowPullTimeout: true,
		Secrets: SecretParameters{
			MaxFileSize:       10485760,
			ExcludedFilenames: nil,
		},
		Registry: RegistryParameters{
			URI:      "",
			Username: "",
			Password: "",
			Token:    "",
		},
		Provenance: "",
	}
}

func IsValidOutputFormat(format string) bool {
	switch OutputFormat(format) {
	case JSON, Table, CycloneDXJSON, CycloneDXXML, SPDXJSON, SPDXXML, SPDXTag, SnapshotJSON:
		return true
	default:
		return false
	}
}
