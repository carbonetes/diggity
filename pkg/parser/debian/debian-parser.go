package debian

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/file"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"

	"github.com/google/uuid"
)

const (
	debType = "deb"
	debian  = "debian"
)

var (
	dpkgStatusPath    = filepath.Join("var", "lib", "dpkg", "status")
	dpkgOldStatusPath = filepath.Join("var", "lib", "dpkg", "status-old")
	dpkgDocPath       = filepath.Join("usr", "share", "doc")
	copyright         = filepath.Join("copyright")
)

// Metadata debian metadata
type Metadata map[string]interface{}

// FindDebianPackagesFromContent Find DPKG packages in the file content
func FindDebianPackagesFromContent() {
	if util.ParserEnabled(debian) {
		for _, content := range file.Contents {
			if strings.Contains(content.Path, dpkgStatusPath) && !strings.Contains(content.Path, dpkgOldStatusPath) {
				if err := readContent(content); err != nil {
					err = errors.New("debian-parser: " + err.Error())
					bom.Errors = append(bom.Errors, &err)
				}
			}
		}
	}
	defer bom.WG.Done()
}

// Read File Contents
func readContent(location *model.Location) error {
	file, err := os.Open(location.Path)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)

	var value string
	var attribute string
	var previousAttribute string

	metadata := make(Metadata)

	for scanner.Scan() {
		keyValue := scanner.Text()

		if strings.Contains(keyValue, ":") {
			keyValues := strings.SplitN(keyValue, ":", 2)
			attribute = keyValues[0]
			value = keyValues[1]

			if strings.Contains(attribute, " ") {
				//clear attribute
				attribute = ""
			}
		} else {
			value = strings.TrimSpace(value + keyValue)
			attribute = previousAttribute
		}

		if len(attribute) > 0 && attribute != " " {
			metadata[attribute] = strings.Replace(value, "\r\n", "", -1)
			metadata[attribute] = strings.Replace(value, "\r ", "", -1)
			metadata[attribute] = strings.TrimSpace(metadata[attribute].(string))
		}

		previousAttribute = attribute

		if len(keyValue) == 0 {
			pkg := new(model.Package)
			pkg.ID = uuid.NewString()
			pkg.Type = debType
			pkg.Path = dpkgStatusPath
			pkg.Locations = append(pkg.Locations, model.Location{
				Path:      util.TrimUntilLayer(*location),
				LayerHash: location.LayerHash,
			})
			// init debian data
			initDebianPackage(pkg, location, metadata)
			bom.Packages = append(bom.Packages, pkg)

			// Reset metadata
			metadata = make(Metadata)
		}

	}

	return nil
}

// Initialize Debian package contents
func initDebianPackage(p *model.Package, location *model.Location, metadata Metadata) *model.Package {

	p.Name = metadata["Package"].(string)
	p.Version = metadata["Version"].(string)
	if val, ok := metadata["Description"].(string); ok {
		p.Description = val
	}

	//check for existing license
	path := strings.Split(location.Path, dpkgStatusPath)[0]
	path = filepath.Join(path, dpkgDocPath, p.Name, copyright)
	searchLicenseOnFileSystem(p, path)
	if len(p.Licenses) > 0 {
		tmpLocation := new(model.Location)
		tmpLocation.LayerHash = location.LayerHash
		tmpLocation.Path = path
		p.Locations = append(p.Locations, model.Location{
			Path:      util.TrimUntilLayer(*tmpLocation),
			LayerHash: location.LayerHash,
		})
	}
	//check files
	if val, ok := metadata["Conffiles"].(string); ok && !*bom.Arguments.DisableFileListing {
		parseDebianFiles(metadata, val)
	}

	//need to add distro in purl
	parseDebianPackageURL(p, metadata["Architecture"].(string))

	//get CPEs
	cpe.NewCPE23(p, "", p.Name, p.Version)

	//fill metadata
	p.Metadata = metadata

	return p
}

// Parse files found on metadata
func parseDebianFiles(m Metadata, filesContent string) {
	lines := strings.Split(filesContent, " ")
	var mapValue = map[string]interface{}{}
	var files []map[string]interface{}
	var path string
	var value string
	var finalValue = map[string]interface{}{}
	for _, line := range lines {
		if strings.Contains(line, "/") {
			path = line
		} else {
			value = line
		}
		if path != "" && value != "" {
			mapValue["value"] = value
			mapValue["algorithm"] = "md5"
			finalValue["path"] = path
			finalValue["digest"] = mapValue

			files = append(files, finalValue)

			// reset map values
			value = ""
			path = ""
			mapValue = map[string]interface{}{}
			finalValue = map[string]interface{}{}
		}
	}
	m["Conffiles"] = files
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

// Parse PURL
func parseDebianPackageURL(pkg *model.Package, architecture string) {
	pkg.PURL = model.PURL("pkg" + ":" + "deb" + "/" + pkg.Name + "@" + pkg.Version + "?arch=" + architecture)
}
