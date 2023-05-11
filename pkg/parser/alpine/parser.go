package alpine

import (
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"

	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	alpine  = "alpine"
	Type = "apk"
)

// Used filepath for path variables
var InstalledPackagesPath = filepath.Join("lib", "apk", "db", "installed")

// Parse installed packages metadata
func parseInstalledPackages(location *model.Location, req *bom.ParserRequirements) error {

	reader, err := os.Open(location.Path)
	if err != nil {
		return err
	}
	defer reader.Close()
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	contents := string(data)
	packages := strings.Split(contents, "\n\n")

	if len(packages) == 0 {
		return nil
	}

	for _, _package := range packages {
		pkg := newPackage(_package, !*req.Arguments.DisableFileListing)
		if pkg == nil {
			continue
		}

		pkg.Locations = append(pkg.Locations, *location)
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

	return nil
}
