package npm

import (
	"path/filepath"

	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const (
	npmPackage         = "package.json"
	npmLock            = "package-lock.json"
	yarnLock           = "yarn.lock"
	invalidPackage     = ".package.json"
	invalidLockPackage = ".package-lock.json"
	invalidYarnlock    = ".yarn.lock"
	Type               = "npm"
)

// FindNpmPackagesFromContent Find NPM packages in the file contents
func FindNpmPackagesFromContent(req *bom.ParserRequirements) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}
	for _, content := range *req.Contents {
		if filepath.Base(content.Path) == npmPackage {
			readNpmContent(&content, req)
		} else if filepath.Base(content.Path) == npmLock {
			readNpmLockContent(&content, req)
		} else if filepath.Base(content.Path) == yarnLock {
			readYarnLockContent(&content, req)
		}
	}

	defer req.WG.Done()
}
