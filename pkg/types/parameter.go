package types

import (
	"fmt"
	"strings"

	"github.com/carbonetes/diggity/internal/helper"
	"github.com/google/go-containerregistry/pkg/name"
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
	ScanType     ScanType     `json:"-"`
	Input        string       `json:"input"`
	OutputFormat OutputFormat `json:"output-format"`
	SaveToFile   string       `json:"save-to-file"`
	Quiet        bool         `json:"quiet"`
	Scanners     []string     `json:"scanners"`
	Provenance   string       `json:"provenance"`
}

func (o OutputFormat) String() string {
	return string(o)
}

func GetAllOutputFormat() string {
	return strings.Join([]string{JSON.String(), Table.String(), CycloneDXJSON.String(), CycloneDXXML.String(), SPDXJSON.String(), SPDXXML.String(), SPDXTag.String(), SnapshotJSON.String()}, ", ")
}

func (p *Parameters) GetScanType() error {
	if _, err := name.ParseReference(p.Input); err == nil {
		p.ScanType = Image
	} else if strings.HasSuffix(p.Input, ".tar") {
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
		return fmt.Errorf("invalid input value %v", p.Input)
	}
	return nil
}

func DefaultParameters() Parameters {
	return Parameters{
		ScanType:     0,
		Input:        "",
		OutputFormat: Table,
		Scanners:     nil,
		Provenance:   "",
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
