package python

import (
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const (
	Type            = "python"
	pythonEgg     = ".egg-info"
	pythonPackage = "METADATA"
	poetry        = "poetry.lock"
	requirements  = "requirements"
)

// FindPythonPackagesFromContent - Find python packages in the file contents
func FindPythonPackagesFromContent(req *bom.ParserRequirements) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}

	for _, content := range *req.Contents {
		if strings.Contains(content.Path, pythonPackage) || strings.Contains(content.Path, pythonEgg) {
			readPythonContent(&content, req)
		}
		if filepath.Base(content.Path) == poetry {
			readPoetryContent(&content, req)
		}
		if strings.Contains(filepath.Base(content.Path), requirements) &&
			strings.Contains(filepath.Base(content.Path), txt) {
			readRequirementsContent(&content, req)
		}
	}

	defer req.WG.Done()
}
