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
func FindAlpinePackagesFromContent(req *bom.ParserRequirements) {
	if util.ParserEnabled(apkType, req.Arguments.EnabledParsers) {
		for _, content := range *req.Contents {
			if strings.Contains(content.Path, installedPackagesPath) {

				if err := parseInstalledPackages(content.Path, content.LayerHash, req.Arguments.DisableFileListing, req.Result.Packages); err != nil {
					err = errors.New("apk-parser: " + err.Error())
					*req.Errors = append(*req.Errors, err)
				}
			}
		}
	}
	defer req.WG.Done()
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
func initAlpinePackage(pkg *model.Package) {
	pkg.Metadata = map[string]string{}
	pkg.ID = uuid.NewString()
	pkg.Type = apkType
	pkg.Path = installedPackagesPath
}

// Parse installed packages metadata
func parseInstalledPackages(filename string, layer string, noFileListing *bool, pkgs *[]model.Package) error {

	var value string
	var attribute string
	var pkg *model.Package = new(model.Package)
	initAlpinePackage(pkg)
	pkg.ID = uuid.NewString()
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
					// pkg.Licenses = strings.Split(value, " ")
					for _, license := range strings.Split(value, " ") {
						if !strings.Contains(strings.ToLower(license), "and") {
							pkg.Licenses = append(pkg.Licenses, license)
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
					pkg.Name = value
				}

			case "S":
				{
					metadata["PackageSize"] = value
				}
			case "T":
				{
					metadata["PackageDescription"] = value
					pkg.Description = value
				}
			case "U":
				{
					metadata["PackageURL"] = value
				}
			case "V":
				{
					metadata["PackageVersion"] = value
					pkg.Version = value
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

			if !*noFileListing {
				files = parseAlpineFiles(apkPackage)
			}

			pkg.Metadata = metadata
		}

		if len(pkg.Metadata.(Manifest)) > 0 {
			if !*noFileListing {
				pkg.Metadata.(Manifest)["Files"] = files
			}

			parseAlpineURL(pkg)

			if pkg.Metadata.(Manifest)["PackageOrigin"] != nil &&
				pkg.Metadata.(Manifest)["PackageName"] != nil &&
				pkg.Metadata.(Manifest)["PackageVersion"] != nil {
				cpe.NewCPE23(pkg,
					pkg.Metadata.(Manifest)["PackageName"].(string),
					pkg.Metadata.(Manifest)["PackageName"].(string),
					pkg.Metadata.(Manifest)["PackageVersion"].(string))

				locations := []model.Location{
					{
						LayerHash: layer,
						Path:      installedPackagesPath,
					},
				}
				for _, content := range file.Contents {
					if strings.Contains(content.Path, pkg.Metadata.(Manifest)["PackageName"].(string)) {
						locations = append(locations, model.Location{
							LayerHash: content.LayerHash,
							Path:      util.TrimUntilLayer(*content),
						})
					}
				}

				pkg.Locations = locations
			}

			// Check if package is not empty before append
			if pkg.Name != "" && pkg.Version != "" {
				*pkgs = append(*pkgs, *pkg)
			}

			files = []model.File{}
			metadata = Manifest{}
			pkg = &model.Package{}
			initAlpinePackage(pkg)
		}

	}
	return nil
}

// Parse PURL
func parseAlpineURL(pkg *model.Package) {
	arch, ok := pkg.Metadata.(Manifest)["Architecture"]
	if !ok {
		arch = ""
	}
	origin, ok := pkg.Metadata.(Manifest)["PackageOrigin"]
	if !ok {
		origin = ""
	}

	pkg.PURL = model.PURL("pkg" + `:` + apkType + `/` + alpine + `/` + pkg.Name + `@` + pkg.Version + `?arch=` + arch.(string) + `&` + `upstream=` + origin.(string) + `&distro=` + alpine)
}
