package alpm

import (
	"path/filepath"

	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const (
	Type      string = "alpm"
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
		if filepath.Base(location.Path) == "desc" {
			parseInstalledPackage(&location, req)
		}
	}
	defer req.WG.Done()
}
