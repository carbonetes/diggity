package maven

import (
	"errors"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/common"
	"github.com/carbonetes/diggity/pkg/parser/util"
	"golang.org/x/exp/maps"
)

const (
	pomFileName      = "pom.xml"
	manifestFile     = "MANIFEST.MF"
	propertiesFile   = "pom.properties"
	Type             = "java"
	jarPackagesRegex = `(?:\.jar)$|(?:\.ear)$|(?:\.war)$|(?:\.jpi)$|(?:\.hpi)$`
)

type (
	// Metadata java metadata
	Metadata map[string]map[string]string
	// Manifest java manifest
	Manifest map[string]string
)

var (
	// JavaPomXML pom metadata
	JavaPomXML          metadata.Project
	nameAndVersionRegex = regexp.MustCompile(`(?Ui)^(?P<name>(?:[[:alpha:]][[:word:].]*(?:\.[[:alpha:]][[:word:].]*)*-?)+)(?:-(?P<version>(?:\d.*|(?:build\d*.*)|(?:rc?\d+(?:^[[:alpha:]].*)?))))?$`)
)

// FindJavaPackagesFromContent checks for jar files in the file contents
func FindJavaPackagesFromContent(req *common.ParserParams) {
	if !util.ParserEnabled(Type, req.Arguments.EnabledParsers) {
		req.WG.Done()
		return
	}
	var result = make(map[string]model.Package, 0)
	for _, content := range *req.Contents {
		if match := regexp.MustCompile(jarPackagesRegex).FindString(content.Path); len(match) > 0 {
			if err := extractJarFile(&content, req.Arguments.Dir, req.DockerTemp, &result); err != nil {
				err = errors.New("java-parser: " + err.Error())
				*req.Errors = append(*req.Errors, err)
			}
		} else if strings.Contains(content.Path, pomFileName) {
			if err := parsePomXML(content, content.Path, req.Arguments.Dir, &result); err != nil {
				err = errors.New("java-parser: " + err.Error())
				*req.Errors = append(*req.Errors, err)
			}
		}
	}
	*req.SBOM.Packages = append(*req.SBOM.Packages, maps.Values(result)...)

	defer req.WG.Done()
}
