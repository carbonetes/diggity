package dpkg

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/common"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

// Read File Contents
func parseDebianPackage(location *model.Location, req *common.ParserParams) {

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
	packages := util.SplitContentsByEmptyLine(contents)
	listFiles := !*req.Arguments.DisableFileListing
	for _, p := range packages {
		metadata := parseMetadata(p, listFiles)
		if metadata == nil {
			continue
		}

		pkg := newPackage(*metadata)
		if pkg == nil {
			continue
		}

		if pkg.Name == "" && pkg.Version == "" {
			continue
		}

		//check for existing license
		path := strings.Split(location.Path, dpkgStatusPath)[0]
		path = filepath.Join(path, dpkgDocPath, pkg.Name, copyright)
		searchLicenseOnFileSystem(pkg, path)
		if len(pkg.Licenses) > 0 {
			tmpLocation := new(model.Location)
			tmpLocation.LayerHash = location.LayerHash
			tmpLocation.Path = path
			pkg.Locations = append(pkg.Locations, model.Location{
				Path:      util.TrimUntilLayer(*tmpLocation),
				LayerHash: location.LayerHash,
			})
		}

		pkg.Path = util.TrimUntilLayer(*location)
		pkg.Locations = append(pkg.Locations, model.Location{
			Path:      pkg.Path,
			LayerHash: location.LayerHash,
		})

		pkg.Metadata = metadata

		*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
	}
}

// Search licenses in existing directory
func searchLicenseOnFileSystem(pkg *model.Package, dpkgDocPath string) {
	// use map license to avoid duplicate entry of license
	var mapLicense = make(map[string]string)
	var value string
	var attribute string
	var licenses []string = make([]string, 0)
	_, err := os.Stat(dpkgDocPath)
	if !os.IsNotExist(err) {
		fileinfo, _ := os.ReadFile(dpkgDocPath)

		lines := strings.Split(string(fileinfo), "\n")
		for _, line := range lines {

			if strings.Contains(line, "License: ") {
				keyValues := strings.Split(line, ": ")
				attribute = keyValues[1]
				value = keyValues[1]
			}

			if len(attribute) > 0 && attribute != " " && value != "none" {
				mapLicense[attribute] = strings.Replace(value, "\r\n", "", -1)
				mapLicense[attribute] = strings.Replace(value, "\r ", "", -1)
				mapLicense[attribute] = strings.TrimSpace(mapLicense[attribute])
			}
		}
		for key := range mapLicense {
			licenses = append(licenses, strings.TrimSpace(key))
		}
	}

	pkg.Licenses = licenses
}
