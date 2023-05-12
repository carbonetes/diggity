package alpm

import (
	"errors"
	"io"
	"os"
	"regexp"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

func parseInstalledPackage(location *model.Location, req *bom.ParserRequirements) {

	file, err := os.Open(location.Path)
	if err != nil {
		err = errors.New("alpmdb-parser: " + err.Error())
		*req.Errors = append(*req.Errors, err)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		err = errors.New("alpmdb-parser: " + err.Error())
		*req.Errors = append(*req.Errors, err)
		return
	}

	contents := string(data)
	attributes := splitAttributes(contents)
	metadata := parseMetadata(attributes)

	if !*req.Arguments.DisableFileListing {
		files := getAlpmFiles(location.Path)
		if files != nil {
			metadata["Files"] = make([]string, len(*files))
			metadata["Files"] = files
		}
	}

	pkg := newPackage(metadata)
	if pkg == nil {
		return
	}

	newLocation := model.Location{
		LayerHash: location.LayerHash,
		Path:      util.TrimUntilLayer(*location),
	}
	pkg.Locations = append(pkg.Locations, newLocation)

	if metadata["License"] != nil {
		pkg.Licenses = append(pkg.Licenses, metadata["License"].([]string)...)
	}

	generateAlpmCpes(pkg)
	*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
}

func splitAttributes(contents string) []string {
	attributes := regexp.
		MustCompile("\r\n").
		ReplaceAllString(contents, "\n")

	return regexp.
		MustCompile(`\n\s*\n`).
		Split(attributes, -1)

}
