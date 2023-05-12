package cargo

import (
	"path/filepath"

	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const (
	Type        string = "rust-crate"
	parserError string = "cargo-parser: "
)

// Metadata cargo metadata
type Metadata map[string]interface{}

// FindCargoPackagesFromContent checks for cargo.lock files in the file contents
func FindCargoPackagesFromContent(req *bom.ParserRequirements) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}

	for _, content := range *req.Contents {
		if filepath.Base(content.Path) == "Cargo.lock" {
			parseCargoPackages(&content, req)
		}
	}

	defer req.WG.Done()
}
