package gradle

import (
	"path/filepath"

	"github.com/carbonetes/diggity/pkg/parser/common"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

var (
	manifestFiles = []string{"buildscript-gradle.lockfile", ".build.gradle"}
	Type          = "gradle"
	Language = "java/kotlin/scala"
)

const parserErr string = "gradle-parser: "

func FindGradlePackagesFromContent(req *common.ParserParams) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}

	for _, content := range *req.Contents {
		base := filepath.Base(content.Path)
		if util.StringSliceContains(manifestFiles, base) {
			parseGradlePackages(&content, req)
		}
	}
	defer req.WG.Done()
}
