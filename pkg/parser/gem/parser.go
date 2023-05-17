package gem

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

func parseGemLockPackage(location *model.Location, req *bom.ParserRequirements) {
	gemFile, err := os.Open(location.Path)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	defer gemFile.Close()

	scanner := bufio.NewScanner(gemFile)
	for scanner.Scan() {
		keyValue := scanner.Text()
		trimedKeyValue := strings.TrimSpace(keyValue)

		if len(keyValue) > 1 && keyValue[0] != ' ' {
			continue
		}

		if isKeyValueValid(keyValue) {
			attributes := strings.Fields(trimedKeyValue)
			if len(attributes) != 2 {
				continue
			}
			pkg := newGemLockPackage(attributes)
			//generate and trim path
			pkg.Locations = append(pkg.Locations, model.Location{
				Path:      util.TrimUntilLayer(*location),
				LayerHash: location.LayerHash,
			})

			*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)

		}
	}
}

// Check if key value is valid
func isKeyValueValid(keyValue string) bool {
	if len(keyValue) < 5 {
		return false
	}
	return strings.Count(keyValue[:5], " ") == 4
}

// Read file contents
func parseGemPackage(location *model.Location, req *bom.ParserRequirements) {
	metadata, err := parseMetadata(location.Path)
	err = errors.New(parserErr + err.Error())
	*req.Errors = append(*req.Errors, err)

	if metadata == nil {
		return
	}

	if len(*metadata) == 0 {
		return
	}

	pkg := newPackage(*metadata)
	if pkg == nil {
		return
	}

	if len(pkg.Name) == 0 || len(pkg.Version) == 0 {
		return
	}

	//generate and trim path
	pkg.Locations = append(pkg.Locations, model.Location{
		Path:      util.TrimUntilLayer(*location),
		LayerHash: location.LayerHash,
	})
	*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)

}
