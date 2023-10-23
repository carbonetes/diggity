package composer

import (
	"path/filepath"

	"github.com/carbonetes/diggity/pkg/parser/common"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const (
	Type     string = "composer"
	Language string = "php"
)

// FindComposerPackagesFromContent - find composers packages from content
func FindComposerPackagesFromContent(req *common.ParserParams) {
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
