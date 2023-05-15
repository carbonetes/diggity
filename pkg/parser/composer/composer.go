package composer

import (
	"strings"

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
		if !strings.Contains(content.Path, composerLock) {
			continue
		}
		parseComposerPackages(&content, req)
	}

	defer req.WG.Done()
}
