package types

import (
	"fmt"
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

type Parameters struct {
	ScanType         ScanType
	Input            string
	OutputFormat     OutputFormat
	SaveToFile       string
	Quiet            bool
	Scanners         []string
	AllowFileListing bool
	MaxFileSize      int64
	// Registry         RegistryParameters `json:"-" yaml:"-"` // ignore this field when marshalling
	Provenance string
}

type RegistryParameters struct {
	// URI      string `json:"uri" yaml:"uri"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	// Token    string `json:"token" yaml:"token"`
}

func (o OutputFormat) String() string {
	return string(o)
}

func GetAllOutputFormat() string {
	return strings.Join([]string{JSON.String(), Table.String(), CycloneDXJSON.String(), CycloneDXXML.String(), SPDXJSON.String(), SPDXXML.String(), SPDXTag.String(), SnapshotJSON.String()}, ", ")
}

func (p *Parameters) GetScanType() error {
	if strings.Contains(p.Input, ":") {
		p.ScanType = Image
	} else if strings.Contains(p.Input, ".tar") {
		p.ScanType = Tarball
	} else if helper.IsDir(p.Input) {
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
		OutputFormat:     Table,
		Quiet:            false,
		Scanners:         nil,
		AllowFileListing: false,
		MaxFileSize:      52428800,
		// Registry: RegistryParameters{
		// 	// URI:      "",
		// 	Username: "",
		// 	Password: "",
		// 	// Token:    "",
		// },
		// Provenance: "",
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
