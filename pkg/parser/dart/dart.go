package dart

import (
	"log"
	"path/filepath"

	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const (
	Type      string = "dart"
	parserErr string = "dart-parser: "
)

// FindDartPackagesFromContent - find dart packages from content
func FindDartPackagesFromContent(req *bom.ParserRequirements) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}

	for _, content := range *req.Contents {
		log.Print(content.Path)
		if filepath.Base(content.Path) == "pubspec.yaml" {
			parseDartPackages(&content, req)
		}

		if filepath.Base(content.Path) == "pubspec.lock" {
			parseDartLockPackages(&content, req)
		}
	}

	defer req.WG.Done()
}
