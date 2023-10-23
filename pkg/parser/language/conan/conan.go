package conan

import (
	"strings"

	"github.com/carbonetes/diggity/pkg/parser/common"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const (
	Type               = "conan"
	parserError string = "conan-parser: "
	Language    string = "c/c++"
)

// FindConanPackagesFromContent Find Conan packages in the file content
func FindConanPackagesFromContent(req *common.ParserParams) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}

	for _, content := range *req.Contents {
		if strings.Contains(content.Path, "conanfile.txt") {
			parseConanFilePackages(&content, req)
		}
		if strings.Contains(content.Path, "conan.lock") {
			parseConanLockPackages(&content, req)
		}
	}

	defer req.WG.Done()
}
