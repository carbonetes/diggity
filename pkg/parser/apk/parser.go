package apk

import (
	"errors"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"

	"io"
	"os"
	"strings"
)

// Parse installed packages metadata
func parseInstalledPackages(location *model.Location, req *bom.ParserRequirements) {

	reader, err := os.Open(location.Path)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}
	defer reader.Close()
	data, err := io.ReadAll(reader)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}
	contents := string(data)
	packages := strings.Split(contents, "\n\n")

	if len(packages) == 0 {
		return
	}

	for _, _package := range packages {
		pkg := newPackage(_package, !*req.Arguments.DisableFileListing)
		if pkg == nil {
			continue
		}
		pkg.Path = util.TrimUntilLayer(*location)
		pkg.Locations = append(pkg.Locations, model.Location{
			Path:      pkg.Path,
			LayerHash: location.LayerHash,
		})
		for _, content := range *req.Contents {
			if strings.Contains(content.Path, pkg.Name) {
				location := model.Location{
					LayerHash: content.LayerHash,
					Path:      util.TrimUntilLayer(content),
				}
				pkg.Locations = append(pkg.Locations, location)
			}
		}
		*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
	}
}
