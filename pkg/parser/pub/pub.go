package pub

import (
	"path/filepath"

	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const (
	Type      string = "pub"
	parserErr string = "pub-parser: "
	Language string = "dart"
)

// FindPubPackagesFromContent - find dart packages from content
func FindPubPackagesFromContent(req *bom.ParserRequirements) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}

	for _, content := range *req.Contents {
		if filepath.Base(content.Path) == "pubspec.yaml" {
			parseDartPackages(&content, req)
		}

		if filepath.Base(content.Path) == "pubspec.lock" {
			parseDartLockPackages(&content, req)
		}
	}

	defer req.WG.Done()
}
