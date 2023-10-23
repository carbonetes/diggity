package nuget

import (
	"strings"

	"github.com/carbonetes/diggity/pkg/parser/common"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const (
	dotnetPackage = ".deps.json"
	Type          = "dotnet"
	nuget         = "nuget"
	Language = "c#/f#/visual_basic"
)

// FindNugetPackagesFromContent - find nuget packages
func FindNugetPackagesFromContent(req *common.ParserParams) {
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
