package dpkg

import (
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/pkg/parser/common"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const (
	Type      string = "deb"
	parserErr string = "dpkg-parser: "
	Distro    string = "debian"
)

var (
	dpkgStatusPath    = filepath.Join("var", "lib", "dpkg", "status")
	dpkgOldStatusPath = filepath.Join("var", "lib", "dpkg", "status-old")
)

// FindDebianPackagesFromContent Find DPKG packages in the file content
func FindDpkgPackagesFromContent(req *common.ParserParams) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}

	for _, content := range *req.Contents {
		if strings.Contains(content.Path, dpkgStatusPath) && !strings.Contains(content.Path, dpkgOldStatusPath) {
			parseDebianPackage(&content, req)
		}
	}

	defer req.WG.Done()
}
