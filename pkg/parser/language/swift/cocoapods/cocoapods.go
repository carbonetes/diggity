package cocoapods

import (
	"path/filepath"

	"github.com/carbonetes/diggity/pkg/parser/common"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const (
	podfilelock = "Podfile.lock"
	pubname     = "name"
	Type        = "pod"
	Language = "swift/objective-c"
)

// FindSwiftPackagesFromContent - find swift and objective-c packages from content
func FindSwiftPackagesFromContent(req *common.ParserParams) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}

	for _, content := range *req.Contents {
		if filepath.Base(content.Path) == podfilelock {
			parseSwiftPackages(&content, req)
		}
	}

	defer req.WG.Done()
}
