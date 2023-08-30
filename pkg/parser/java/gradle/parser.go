package gradle

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

func parseGradlePackages(location *model.Location, req *bom.ParserRequirements) {
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

		if !strings.Contains(properties, ":") {
			continue
		}

		var gradleMetadata metadata.GradleMetadata
		values := strings.SplitN(properties, ":", 3)

		if len(values) != 3 {
			continue
		}

		gradleMetadata.Vendor = values[0]
		gradleMetadata.Name = values[1]
		gradleMetadata.Version = strings.ReplaceAll(values[2], "=classpath", "")

		pkg := newPackage(gradleMetadata)
		pkg.Path = util.TrimUntilLayer(*location)
		pkg.Locations = append(pkg.Locations, model.Location{
			Path: pkg.Path,
			LayerHash: location.LayerHash,
		})
		*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
	}

}