package gobin

import (
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/golang"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

// FindGoBinPackagesFromContent Find go binaries in the file contents
func FindGoBinPackagesFromContent(req *bom.ParserRequirements) {
	if !util.ParserEnabled(golang.Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}
	// Look for go bin file
	for _, content := range *req.Contents {
		if strings.Contains(filepath.Base(content.Path), "gobin") {
			parseGoBinContent(&content, req)
		}
	}

	defer req.WG.Done()
}
