package gomod

import (
	"path/filepath"

	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/golang"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

// FindGoModPackagesFromContent Find go.mod in the file contents
func FindGoModPackagesFromContent(req *bom.ParserRequirements) {
	if !util.ParserEnabled(golang.Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}

	for _, content := range *req.Contents {
		if filepath.Base(content.Path) == "go.mod" {
			parseGoModContent(&content, req)
		}
	}

	defer req.WG.Done()
}
