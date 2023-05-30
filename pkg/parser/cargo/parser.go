package cargo

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

func parseCargoPackages(location *model.Location, req *bom.ParserRequirements) {
	file, err := os.Open(location.Path)
	if err != nil {
		err = errors.New(parserError + err.Error())
		*req.Errors = append(*req.Errors, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		err = errors.New(parserError + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	contents := string(data)
	contents = strings.NewReplacer(`"`, ``, `,`, ``, ` `, ``).Replace(contents)
	packages := util.SplitContentsByEmptyLine(contents)
	for _, p := range packages {
		if !strings.Contains(p, "[[package]]") {
			continue
		}
		metadata := parseMetadata(p)
		if metadata == nil {
			continue
		}
		pkg := newPackage(location, *metadata)
		if pkg == nil {
			continue
		}
		*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
	}
}
