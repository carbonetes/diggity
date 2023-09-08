package sbt

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

func parserSbtPackages(location *model.Location, req *bom.ParserRequirements) {
	file, err := os.Open(location.Path)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
		return
	}

	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		properties := scanner.Text()

		if !strings.Contains(properties, "%") {
			continue
		}

		var sbtMetadata metadata.SbtMetadata
		values := strings.SplitN(properties, "%", 3)
		
		if len(values) != 3 && len(values) <= 4 {
			continue
		}

		if strings.Contains(values[2], "%") {
			splitNameVersion := strings.SplitN(values[2], "%", 2)
			values[2] = splitNameVersion[1]
			if len(values[1]) == 0 {
				values[1] = splitNameVersion[0]
			}
		}

		sbtMetadata.Vendor = removeDoubleQoute(values[0])
		sbtMetadata.Name = removeDoubleQoute(values[1])
		sbtMetadata.Version = removeDoubleQoute(values[2])
		if strings.Contains(sbtMetadata.Version, "%") {
			splitVersionConfig := strings.SplitN(sbtMetadata.Version, "%", 2)
			sbtMetadata.Version = splitVersionConfig[0]
			sbtMetadata.Config = splitVersionConfig[1]
		}

		sbtMetadata.Version = sanitizeVersion(sbtMetadata.Version)

		pkg := newPackage(sbtMetadata)
		pkg.Path = util.TrimUntilLayer(*location)
		pkg.Locations = append(pkg.Locations, model.Location{
			Path:      pkg.Path,
			LayerHash: location.LayerHash,
		})
		*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
		
	}
}

func removeDoubleQoute(word string) string{
	word = strings.ReplaceAll(word, "\"", "")
	return strings.ReplaceAll(word, ",", "")
}

func sanitizeVersion(version string) string {
	// Define a regular expression to match allowed characters for special cases sbt dependencies
	allowedPattern := regexp.MustCompile(`[^a-zA-Z0-9._-]+`) // allowed ^._-, all numbers and letters
	return allowedPattern.ReplaceAllString(version, "")
}