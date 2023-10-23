package hex

import (
	"errors"
	"path/filepath"

	"github.com/carbonetes/diggity/pkg/parser/common"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const (
	rebarLock = "rebar.lock"
	mixLock   = "mix.lock"
	Type      = "hex"
	Language  = "erlang"
)

// FindHexPackagesFromContent - find hex packages from content
func FindHexPackagesFromContent(req *common.ParserParams) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}

	for _, content := range *req.Contents {
		if filepath.Base(content.Path) == rebarLock {
			if err := parseHexRebarPacakges(&content, req.SBOM.Packages); err != nil {
				err = errors.New("hex-parser: " + err.Error())
				*req.Errors = append(*req.Errors, err)
			}
		}
		if filepath.Base(content.Path) == mixLock {
			if err := parseHexMixPackages(&content, req.SBOM.Packages); err != nil {
				err = errors.New("hex-parser: " + err.Error())
				*req.Errors = append(*req.Errors, err)
			}
		}
	}

	defer req.WG.Done()
}
