package pacman

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

func parseInstalledPackage(location *model.Location, req *bom.ParserRequirements) {

	file, err := os.Open(location.Path)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
		return
	}

	contents := string(data)
	attributes := util.SplitContentsByEmptyLine(contents)
	metadata := parseMetadata(attributes)

	if !*req.Arguments.DisableFileListing {
		files, backups, err := getAlpmFiles(location.Path)
		if err != nil {
			err = errors.New(parserErr + err.Error())
			*req.Errors = append(*req.Errors, err)
		}

		if files != nil && len(*files) > 0 {
			metadata["Files"] = make([]string, len(*files))
			metadata["Files"] = files
		}
		if backups != nil && len(*backups) > 0 {
			metadata["Backups"] = make([]string, len(*files))
			metadata["Backups"] = backups
		}
	}

	pkg := newPackage(metadata)
	if pkg == nil {
		return
	}

	pkg.Path = strings.Replace(util.TrimUntilLayer(*location), "\\desc", "", -1)

	newLocation := model.Location{
		LayerHash: location.LayerHash,
		Path:      util.TrimUntilLayer(*location),
	}
	pkg.Locations = append(pkg.Locations, newLocation)

	if metadata["Licenses"] != nil {
		pkg.Licenses = append(pkg.Licenses, metadata["Licenses"].([]string)...)
	}

	generatePacmanCpes(pkg)
	*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
}
