package hackage

import (
	"path/filepath"

	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const (
	Type       = "hackage"
	stackYaml     = "stack.yaml"
	stackYamlLock = "stack.yaml.lock"
	cabalFreeze   = "cabal.project.freeze"
	shaTag        = "sha256"
	revTag        = "rev"
	anyTag        = "any."
	constraints   = "constraints:"
)

// FindHackagePackagesFromContent checks for stack.yaml, stack.yaml.lock, and cabal.project.freeze files in the file contents
func FindHackagePackagesFromContent(req *bom.ParserRequirements) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}

	for _, content := range *req.Contents {
		if filepath.Base(content.Path) == stackYaml {
			readStackContent(&content, req)
		}
		if filepath.Base(content.Path) == stackYamlLock {
			readStackLockContent(&content, req)
		}
		if filepath.Base(content.Path) == cabalFreeze {
			readCabalFreezeContent(&content, req)
		}
	}

	defer req.WG.Done()
}
