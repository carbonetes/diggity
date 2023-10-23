package hackage

import (
	"path/filepath"

	"github.com/carbonetes/diggity/pkg/parser/common"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const (
	Type          string = "hackage"
	stackYaml     string = "stack.yaml"
	stackYamlLock string = "stack.yaml.lock"
	cabalFreeze   string = "cabal.project.freeze"
	shaTag        string = "sha256"
	revTag        string = "rev"
	anyTag        string = "any."
	constraints   string = "constraints:"
	Language      string = "haskell"
)

// FindHackagePackagesFromContent checks for stack.yaml, stack.yaml.lock, and cabal.project.freeze files in the file contents
func FindHackagePackagesFromContent(req *common.ParserParams) {
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
