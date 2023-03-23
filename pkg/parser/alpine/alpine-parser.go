package alpine

import (
	"errors"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/file"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"

	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const (
	alpine  = "alpine"
	apkType = "apk"
)

// Used filepath for path variables
var installedPackagesPath = filepath.Join("lib", "apk", "db", "installed")

// Manifest alpine manifest
type Manifest map[string]interface{}

// FindAlpinePackagesFromContent check for alpine-os files in the file contents
func FindAlpinePackagesFromContent() {
	if util.ParserEnabled(apkType) {
		for _, content := range file.Contents {
			if strings.Contains(content.Path, installedPackagesPath) {
				if err := parseInstalledPackages(content.Path, content.LayerHash); err != nil {
					err = errors.New("apk-parser: " + err.Error())
					bom.Errors = append(bom.Errors, &err)
				}
			}
		}
	}
	defer bom.WG.Done()
}

// Parse alpine files
func parseAlpineFiles(content string) []model.File {

	var files []model.File
	keyValues := strings.Split(content, "\n")
	for idx := range keyValues {
		file := model.File{}

		// F: = File or Directory
		if strings.HasPrefix(keyValues[idx], "F:") {
			file.Path = strings.SplitN(keyValues[idx], ":", 2)[1]
			files = append(files, file)
		} else if strings.HasPrefix(keyValues[idx], "R:") {
			// reloop until F: or R: prefix is found
			file.Path = strings.SplitN(keyValues[idx], ":", 2)[1]
			for fileIdx := idx + 1; fileIdx < len(keyValues) && !strings.HasPrefix(keyValues[fileIdx], "R:"); {
				//  a:, M: = File Permissions
				if strings.HasPrefix(keyValues[fileIdx], "a:") || strings.HasPrefix(keyValues[fileIdx], "M:") {
					file.OwnerGID = strings.Split(keyValues[fileIdx], ":")[1]
					file.OwnerUID = strings.Split(keyValues[fileIdx], ":")[2]
					file.Permissions = strings.Split(keyValues[fileIdx], ":")[3]
				} else if /* Z: = Pull Checksum */ strings.HasPrefix(keyValues[fileIdx], "Z:") {
					digest := map[string]string{}
					digest["algorithm"] = "sha1"
					digest["value"] = strings.SplitN(keyValues[fileIdx], ":", 2)[1]
					file.Digest = digest
				}
				fileIdx++
			}
			files = append(files, file)
		}

	}

	return files
}

// Init alpine package
func initAlpinePackage(_package *model.Package) {
	_package.Metadata = map[string]string{}
	_package.ID = uuid.NewString()
	_package.Type = apkType
	_package.Path = installedPackagesPath
}

// Parse installed packages metadata
func parseInstalledPackages(filename string, layer string) error {

	var value string
	var attribute string
	var _package *model.Package = new(model.Package)
	initAlpinePackage(_package)
	_package.ID = uuid.NewString()
	var files []model.File
	metadata := Manifest{}

	reader, err := os.Open(filename)
	if err != nil {
		return err
	}
	bContent, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	sContent := string(bContent)
	apkPackages := strings.Split(sContent, "\n\n")
	for _, apkPackage := range apkPackages {
		apkPackage = strings.TrimSpace(apkPackage)
		keyValues := strings.Split(apkPackage, "\n")
		for _, keyValue := range keyValues {

			if strings.Contains(keyValue, ":") && !strings.Contains(keyValue, ":=") {
				keyValues := strings.SplitN(keyValue, ":", 2)
				attribute = keyValues[0]
				value = keyValues[1]
			} else {
				value = strings.TrimSpace(value + keyValue)
			}

			switch attribute {
			case "A":
				{
					metadata["Architecture"] = value
				}
			case "C":
				{
					metadata["PullChecksum"] = value
				}
			case "D", "r":
				{
					metadata["PullDependencies"] = value
				}
			case "I":
				{
					metadata["PackageInstalledSize"] = value
				}
			case "L":
				{
					metadata["License"] = value
					// _package.Licenses = strings.Split(value, " ")
					for _, license := range strings.Split(value, " ") {
						if !strings.Contains(strings.ToLower(license), "and") {
							_package.Licenses = append(_package.Licenses, license)
						}
					}
				}
			case "M":
				{
					metadata["Permissions"] = value
				}
			case "P":
				{
					metadata["PackageName"] = value
					_package.Name = value
				}

			case "S":
				{
					metadata["PackageSize"] = value
				}
			case "T":
				{
					metadata["PackageDescription"] = value
					_package.Description = value
				}
			case "U":
				{
					metadata["PackageURL"] = value
				}
			case "V":
				{
					metadata["PackageVersion"] = value
					_package.Version = value
				}
			case "c":
				{
					metadata["GitCommitHashApk"] = value
				}
			case "m":
				{
					metadata["Maintainer"] = value
				}
			case "o":
				{
					metadata["PackageOrigin"] = value
				}
			case "p":
				{
					metadata["Provides"] = value
				}
			case "t":
				{
					metadata["BuildTimestamp"] = value
				}
			}

			if !*bom.Arguments.DisableFileListing {
				files = parseAlpineFiles(apkPackage)
			}

			_package.Metadata = metadata
		}

		if len(_package.Metadata.(Manifest)) > 0 {
			if !*bom.Arguments.DisableFileListing {
				_package.Metadata.(Manifest)["Files"] = files
			}

			parseAlpineURL(_package)

			if _package.Metadata.(Manifest)["PackageOrigin"] != nil &&
				_package.Metadata.(Manifest)["PackageName"] != nil &&
				_package.Metadata.(Manifest)["PackageVersion"] != nil {
				cpe.NewCPE23(_package,
					_package.Metadata.(Manifest)["PackageName"].(string),
					_package.Metadata.(Manifest)["PackageName"].(string),
					_package.Metadata.(Manifest)["PackageVersion"].(string))

				locations := []model.Location{
					{
						LayerHash: layer,
						Path:      installedPackagesPath,
					},
				}
				for _, content := range file.Contents {
					if strings.Contains(content.Path, _package.Metadata.(Manifest)["PackageName"].(string)) {
						locations = append(locations, model.Location{
							LayerHash: content.LayerHash,
							Path:      util.TrimUntilLayer(*content),
						})
					}
				}

				_package.Locations = locations
			}

			// Check if package is not empty before append
			if _package.Name != "" && _package.Version != "" {
				bom.Packages = append(bom.Packages, _package)
			}

			files = []model.File{}
			metadata = Manifest{}
			_package = &model.Package{}
			initAlpinePackage(_package)
		}

	}
	return nil
}

// Parse PURL
func parseAlpineURL(_package *model.Package) {
	arch, ok := _package.Metadata.(Manifest)["Architecture"]
	if !ok {
		arch = ""
	}
	origin, ok := _package.Metadata.(Manifest)["PackageOrigin"]
	if !ok {
		origin = ""
	}

	_package.PURL = model.PURL("pkg" + `:` + apkType + `/` + alpine + `/` + _package.Name + `@` + _package.Version + `?arch=` + arch.(string) + `&` + `upstream=` + origin.(string) + `&distro=` + alpine)
}
