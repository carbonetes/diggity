package alpine

import (
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const Type string = "apk"

// Used filepath for path variables
var InstalledPackagesPath = filepath.Join("lib", "apk", "db", "installed")

// FindAlpinePackagesFromContent check for alpine-os files in the file contents
func FindAlpinePackagesFromContent(req *bom.ParserRequirements) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}

	for _, content := range *req.Contents {
		if strings.Contains(content.Path, InstalledPackagesPath) {
			parseInstalledPackages(&content, req)
		}
	}
	defer req.WG.Done()
}
