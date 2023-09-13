package nix

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
	"golang.org/x/exp/slices"
)

const Type = "nix"

var nixStorePath = filepath.Join("nix", "store")

func FindNixPackagesFromContent(req *bom.ParserRequirements) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}
	var paths []string
	for _, content := range *req.Contents {
		if strings.Contains(content.Path, nixStorePath) {
			target := getPackageDir(content.Path)
			if len(target) == 0 || slices.Contains(paths, target) {
				continue
			}
			paths = append(paths, target)
			parseNixPackage(target, &content, req)
		}
	}

	defer req.WG.Done()
}

func getPackageDir(content string) string {
	if strings.Contains(filepath.Base(content), ".nix") || strings.Contains(filepath.Base(content), ".drv") {
		return ""
	}
	fields := strings.Split(content, string(os.PathSeparator))
	for index, field := range fields {
		if field == "nix" {
			if fields[index+1] == "store" && index+2 < len(fields) {
				return fields[index+2]
			}

		}
	}

	return ""
}
