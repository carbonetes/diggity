package alpine

import (
	"errors"
	"strings"

	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

// FindAlpinePackagesFromContent check for alpine-os files in the file contents
func FindAlpinePackagesFromContent(req *bom.ParserRequirements) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		return
	}
	for _, content := range *req.Contents {
		if !strings.Contains(content.Path, InstalledPackagesPath) {
			continue
		}
		if err := parseInstalledPackages(&content, req); err != nil {
			err = errors.New("apk-parser: " + err.Error())
			*req.Errors = append(*req.Errors, err)
		}

	}

	defer req.WG.Done()
}
