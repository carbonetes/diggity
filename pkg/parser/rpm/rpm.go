package rpm

import (
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const Type = "rpm"

var (
	rpmPackagesPath = filepath.Join("rpm", "Packages")
)

// FindRpmPackagesFromContent Find rpm/Packages in the file content.
func FindRpmPackagesFromContent(req *bom.ParserRequirements) {
	// Get RPM Information if rpm/Packages is found
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}

	for _, content := range *req.Contents {
		if strings.Contains(content.Path, rpmPackagesPath) {
			readRpmContent(&content, req)
		}
	}

	defer req.WG.Done()
}
