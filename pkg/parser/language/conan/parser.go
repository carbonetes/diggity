package conan

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/common"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

func parseConanFilePackages(location *model.Location, req *common.ParserParams) {
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
	attributes := util.SplitContentsByEmptyLine(contents)
	metadataList := new([]metadata.ConanMetadata)
	for _, attribute := range attributes {
		if !strings.Contains(attribute, "[requires]") {
			continue
		}
		metadataList = parseConanFileMetadata(attribute[1:])
	}

	if metadataList == nil {
		return
	}

	if len(*metadataList) == 0 {
		return
	}

	for _, metadata := range *metadataList {
		pkg := newPackage(location, metadata)
		if pkg == nil {
			continue
		}

		if len(pkg.Name) == 0 {
			continue
		}

		*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
	}
}

func parseConanLockPackages(location *model.Location, req *common.ParserParams) {

	metadata, err := parseConanLockMetadata(location.Path)
	if err != nil {
		err = errors.New(parserError + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	if metadata == nil {
		return
	}

	if len(metadata.GraphLock.Nodes) == 0 {
		return
	}

	for _, conanPkg := range metadata.GraphLock.Nodes {
		if conanPkg.Ref == "" {
			continue
		}

		pkg := newPackage(location, conanPkg)

		if pkg == nil {
			continue
		}

		if len(pkg.Name) == 0 {
			continue
		}

		*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
	}

}
