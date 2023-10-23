package rpm

import (
	"path/filepath"

	"github.com/carbonetes/diggity/pkg/parser/common"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const Type = "rpm"

var (
	rpmDBFiles = []string{"Packages", "Packages.db", "rpmdb.sqlite"}
)

// FindRpmPackagesFromContent Find rpm/Packages in the file content.
func FindRpmPackagesFromContent(req *common.ParserParams) {
	// Get RPM Information if rpm/Packages is found
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}

	for _, content := range *req.Contents {
		base := filepath.Base(content.Path)
		if util.StringSliceContains(rpmDBFiles, base) {
			readRpmContent(&content, req)
		}
	}

	defer req.WG.Done()
}
