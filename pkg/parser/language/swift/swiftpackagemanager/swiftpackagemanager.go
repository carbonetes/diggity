package swiftpackagemanager

import (
	"path/filepath"

	"github.com/carbonetes/diggity/pkg/parser/common"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

var (
	manifestFiles = []string{"Package.resolved", ".package.resolved"}
	Type          = "swift"
)

const parserErr string = "swift-parser: "

func FindSwiftPackagesFromContent(req *common.ParserParams) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}

	for _, content := range *req.Contents {
		base := filepath.Base(content.Path)
		if util.StringSliceContains(manifestFiles, base) {
			parseSwiftPackages(&content, req)
		}
	}
	defer req.WG.Done()
}
