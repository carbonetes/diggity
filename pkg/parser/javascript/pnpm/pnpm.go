package pnpm

import (
	"path/filepath"

	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

var (
	pnpmLockYamlFile = "pnpm-lock.yaml"
	Type             = "pnpm"
)

func FindPnpmPackagesFromContent(req *bom.ParserRequirements) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}
	for _, content := range *req.Contents {
		if filepath.Base(content.Path) == pnpmLockYamlFile {
			readLockFile(&content, req)
		}
	}

	defer req.WG.Done()
}
