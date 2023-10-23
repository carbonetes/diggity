package portage

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/pkg/parser/common"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

var (
	portageDBPath = filepath.Join("var", "db", "pkg") + string(os.PathSeparator)
	Type          = "portage"
)

// FindPortagePackagesFromContent find portage metadata files
func FindPortagePackagesFromContent(req *common.ParserParams) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}

	for _, content := range *req.Contents {
		if strings.Contains(content.Path, portageDBPath) &&
			strings.Contains(content.Path, portageContent) {
			readPortageContent(&content, req)
		}
	}

	defer req.WG.Done()
}
