package sbt

import (
	"path/filepath"

	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

// simple build tool - scala
var (
	manifestFiles = []string{"build.sbt"}
	Type = "sbt"
)

const parserErr string = "sbt-parser: "

func FindSbtPackagesFromContent(req *bom.ParserRequirements) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}

	for _, content := range *req.Contents {
		base := filepath.Base(content.Path)
		if util.StringSliceContains(manifestFiles, base) {
			parserSbtPackages(&content, req)
		}
	}
	defer req.WG.Done()
}