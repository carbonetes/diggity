package gem

import (
	"strings"

	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const (
	Type      string = "gem"
	parserErr string = "gem-parser: "
)

// FindGemPackagesFromContent Find gem packages in the file contents
func FindGemPackagesFromContent(req *bom.ParserRequirements) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}

	for _, content := range *req.Contents {
		if strings.Contains(content.Path, ".gemspec") && strings.Contains(content.Path, "specifications") {
			parseGemPackage(&content, req)
		}
		if strings.Contains(content.Path, "Gemfile.lock") {
			parseGemLockPackage(&content, req)
		}
	}
	defer req.WG.Done()
}
