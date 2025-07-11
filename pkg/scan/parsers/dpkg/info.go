package dpkg

import (
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// InfoParser handles parsing of /var/lib/dpkg/info/* files
type InfoParser struct {
	base *common.BaseParser
}

// NewInfoParser creates a new DPKG info parser
func NewInfoParser(base *common.BaseParser) *InfoParser {
	return &InfoParser{base: base}
}

// Parse parses DPKG info files (.list, .md5sums, .conffiles, etc.)
func (p *InfoParser) Parse(filePath string) ([]model.Package, error) {
	// TODO: Implement DPKG info parsing
	// These files contain additional package information like file lists
	return []model.Package{}, nil
}
