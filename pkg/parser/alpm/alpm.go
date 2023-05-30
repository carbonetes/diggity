package alpm

import (
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const (
	Type   string = "alpm"
	parserErr string = "alpm-parser: "
)

var InstalledPackagesPath = filepath.Join("var", "lib", "pacman", "local")

type Metadata map[string]interface{}

func FindAlpmPackagesFromContent(req *bom.ParserRequirements) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}

	for _, location := range *req.Contents {
		if strings.Contains(location.Path, InstalledPackagesPath) && strings.Contains(location.Path, "\\desc") {
			parseInstalledPackage(&location, req)
		}
	}
	defer req.WG.Done()
}
