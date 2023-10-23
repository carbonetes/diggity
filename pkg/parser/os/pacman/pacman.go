package pacman

import (
	"path/filepath"

	"github.com/carbonetes/diggity/pkg/parser/common"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const (
	Type      string = "pacman"
	parserErr string = "pacman-parser: "
)

var InstalledPackagesPath = filepath.Join("var", "lib", "pacman", "local")

type Metadata map[string]interface{}

func FindPacmanPackagesFromContent(req *common.ParserParams) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}

	for _, location := range *req.Contents {
		if filepath.Base(location.Path) == "desc" {
			parseInstalledPackage(&location, req)
		}
	}
	defer req.WG.Done()
}
