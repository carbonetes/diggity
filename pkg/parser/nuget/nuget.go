package nuget

import (
	"strings"

	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const (
	dotnetPackage = ".deps.json"
	Type          = "dotnet"
	nuget         = "nuget"
)

// FindNugetPackagesFromContent - find nuget packages
func FindNugetPackagesFromContent(req *bom.ParserRequirements) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}

	for _, content := range *req.Contents {
		if strings.Contains(content.Path, dotnetPackage) {
			parseNugetPackages(&content, req)
		}
	}

	defer req.WG.Done()
}
