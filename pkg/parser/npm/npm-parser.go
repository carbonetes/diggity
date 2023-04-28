package npm

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"

	"github.com/google/uuid"
)

const (
	npmPackage         = "package.json"
	npmLock            = "package-lock.json"
	yarnLock           = "yarn.lock"
	invalidPackage     = ".package.json"
	invalidLockPackage = ".package-lock.json"
	invalidYarnlock    = ".yarn.lock"
	npm                = "npm"
)

// LockMetadata npm lock metadata type
type LockMetadata map[string]interface{}

var (
	// NpmMetadata  metadata
	NpmMetadata metadata.PackageJSON
	// NpmLockMetadata lock metadata
	NpmLockMetadata metadata.PackageLock
	packageRegEx    = regexp.MustCompile(`^"?((?:@\w[\w-_.]*\/)?\w[\w-_.]*)@`)
	versionRegEx    = regexp.MustCompile(`^\W+version(?:\W+"|:\W+)([\w-_.]+)"?`)
)

// FindNpmPackagesFromContent Find NPM packages in the file contents
func FindNpmPackagesFromContent(req *bom.ParserRequirements) {
	if util.ParserEnabled(npm, req.Arguments.EnabledParsers) {
		for _, content := range *req.Contents {
			if filepath.Base(content.Path) == npmPackage {
				if err := readNpmContent(&content, req.SBOM.Packages); err != nil {
					err = errors.New("npm-parser: " + err.Error())
					*req.Errors = append(*req.Errors, err)
				}
			} else if filepath.Base(content.Path) == npmLock {
				if err := readNpmLockContent(&content, req.SBOM.Packages); err != nil {
					err = errors.New("npm-parser: " + err.Error())
					*req.Errors = append(*req.Errors, err)
				}
			} else if filepath.Base(content.Path) == yarnLock {
				if err := readYarnLockContent(&content, req.SBOM.Packages); err != nil {
					err = errors.New("npm-parser: " + err.Error())
					*req.Errors = append(*req.Errors, err)
				}
			}
		}
	}
	defer req.WG.Done()
}

// Read file contents
func readNpmContent(location *model.Location, pkgs *[]model.Package) error {
	file, err := os.ReadFile(location.Path)
	if err != nil {
		return err
	}

	// Validate if file is valid JSON
	if !json.Valid(file) {
		return nil
	}

	if err = json.Unmarshal(file, &NpmMetadata); err != nil {
		return err
	}

	if NpmMetadata.Name != "" {
		pkg := new(model.Package)
		pkg.ID = uuid.NewString()
		pkg.Locations = append(pkg.Locations, model.Location{
			Path:      util.TrimUntilLayer(*location),
			LayerHash: location.LayerHash,
		})

		// // init npm data
		pkg.Name = NpmMetadata.Name
		pkg.Version = NpmMetadata.Version
		pkg.Description = NpmMetadata.Description
		pkg.Type = npm
		pkg.Path = NpmMetadata.Name

		// // check type of license then parse
		switch NpmMetadata.License.(type) {
		case string:
			pkg.Licenses = append(pkg.Licenses, NpmMetadata.License.(string))
		case map[string]interface{}:
			license := NpmMetadata.License.(map[string]interface{})
			if _, ok := license["type"]; ok {
				pkg.Licenses = append(pkg.Licenses, license["type"].(string))
			}
		}

		// //parseURL
		parseNpmPackageURL(pkg)
		cpe.NewCPE23(pkg, pkg.Name, pkg.Name, pkg.Version)
		pkg.Metadata = NpmMetadata

		*pkgs = append(*pkgs, *pkg)

	}
	return nil
}

// Parse lock content
func readNpmLockContent(location *model.Location, pkgs *[]model.Package) error {

	file, err := os.ReadFile(location.Path)
	if err != nil {
		return err
	}

	// Validate if file is valid JSON
	if !json.Valid(file) {
		return nil
	}

	if err = json.Unmarshal(file, &NpmLockMetadata); err != nil {
		return err
	}

	if len(NpmLockMetadata.Dependencies) > 0 {
		for name, cPackage := range NpmLockMetadata.Dependencies {
			pkg := new(model.Package)
			pkg.ID = uuid.NewString()
			pkg.Locations = append(pkg.Locations, model.Location{
				Path:      util.TrimUntilLayer(*location),
				LayerHash: location.LayerHash,
			})

			// // init npm data
			pkg.Name = name
			pkg.Version = cPackage.Version
			pkg.Type = npm
			pkg.Path = name

			// //parseURL
			parseNpmPackageURL(pkg)
			cpe.NewCPE23(pkg, pkg.Name, pkg.Name, pkg.Version)
			pkg.Metadata = cPackage

			*pkgs = append(*pkgs, *pkg)

		}
	}

	return nil
}

// Parse yarn lock content
func readYarnLockContent(location *model.Location, pkgs *[]model.Package) error {

	file, err := os.Open(location.Path)
	if err != nil {
		return err
	}
	metadata := make(LockMetadata)
	scanner := bufio.NewScanner(file)

	var value string
	var attribute string

	for scanner.Scan() {
		keyValue := scanner.Text()

		packageMatches := packageRegEx.FindStringSubmatch(keyValue)
		if len(packageMatches) >= 2 {
			attribute = "Name"
			value = packageMatches[1]
		}

		versioMatches := versionRegEx.FindStringSubmatch(keyValue)
		if len(versioMatches) >= 2 {
			attribute = "Version"
			value = versioMatches[1]
		}

		if len(attribute) > 0 && attribute != " " {
			metadata[attribute] = strings.Replace(value, "\r\n", "", -1)
			metadata[attribute] = strings.Replace(value, "\r ", "", -1)
			metadata[attribute] = strings.TrimSpace(metadata[attribute].(string))
		}

		if _, ok := metadata["Name"].(string); ok && len(keyValue) == 0 && len(metadata) >= 2 {
			pkg := new(model.Package)
			pkg.ID = uuid.NewString()
			pkg.Type = npm
			pkg.Name = metadata["Name"].(string)
			pkg.Path = metadata["Name"].(string)

			if metadata["Version"] != nil {
				pkg.Version = metadata["Version"].(string)
			}

			pkg.Locations = append(pkg.Locations, model.Location{
				Path:      util.TrimUntilLayer(*location),
				LayerHash: location.LayerHash,
			})
			parseNpmPackageURL(pkg)
			cpe.NewCPE23(pkg, pkg.Name, pkg.Name, pkg.Version)
			pkg.Metadata = metadata
			*pkgs = append(*pkgs, *pkg)
			metadata = LockMetadata{}
		}
	}

	return nil
}

// Parse PURL
func parseNpmPackageURL(pkg *model.Package) {
	pkg.PURL = model.PURL("pkg" + ":" + npm + "/" + pkg.Name + "@" + pkg.Version)
}
