package composer

import (
	"path/filepath"

	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const Type = "php"

// FindComposerPackagesFromContent - find composers packages from content
func FindComposerPackagesFromContent(req *bom.ParserRequirements) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}

	for _, content := range *req.Contents {
		if filepath.Base(content.Path) == composerLock {
			parseComposerPackages(&content, req)
		}
	}

	defer req.WG.Done()
}
