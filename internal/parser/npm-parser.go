package parser

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/file"
	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/model/metadata"

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

// NpmLockMetadata npm lock metadata type
type NpmLockMetadata map[string]interface{}

var (
	// NpmMetadata  metadata
	NpmMetadata metadata.PackageJSON
	// NPMLockMetadata lock metadata
	NPMLockMetadata metadata.PackageLock
	packageRegEx    = regexp.MustCompile(`^"?((?:@\w[\w-_.]*\/)?\w[\w-_.]*)@`)
	versionRegEx    = regexp.MustCompile(`^\W+version(?:\W+"|:\W+)([\w-_.]+)"?`)
)

// FindNpmPackagesFromContent Find DPKG packages in the file contents
func FindNpmPackagesFromContent() {
	if parserEnabled(npm) {
		for _, content := range file.Contents {
			if filepath.Base(content.Path) == npmPackage {
				if err := readNpmContent(content); err != nil {
					err = errors.New("npm-parser: " + err.Error())
					Errors = append(Errors, &err)
				}
			} else if filepath.Base(content.Path) == npmLock {
				if err := readNpmLockContent(content); err != nil {
					err = errors.New("npm-parser: " + err.Error())
					Errors = append(Errors, &err)
				}
			} else if filepath.Base(content.Path) == yarnLock {
				if err := readYarnLockContent(content); err != nil {
					err = errors.New("npm-parser: " + err.Error())
					Errors = append(Errors, &err)
				}
			}
		}
	}
	defer WG.Done()
}

// Read file contents
func readNpmContent(location *model.Location) error {
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
		_package := new(model.Package)
		_package.ID = uuid.NewString()
		_package.Locations = append(_package.Locations, model.Location{
			Path:      TrimUntilLayer(*location),
			LayerHash: location.LayerHash,
		})

		// // init npm data
		_package.Name = NpmMetadata.Name
		_package.Version = NpmMetadata.Version
		_package.Description = NpmMetadata.Description
		_package.Type = npm
		_package.Path = NpmMetadata.Name

		// // check type of license then parse
		switch NpmMetadata.License.(type) {
		case string:
			_package.Licenses = append(_package.Licenses, NpmMetadata.License.(string))
		case map[string]interface{}:
			license := NpmMetadata.License.(map[string]interface{})
			if _, ok := license["type"]; ok {
				_package.Licenses = append(_package.Licenses, license["type"].(string))
			}
		}

		// //parseURL
		parseNpmPackageURL(_package)
		cpe.NewCPE23(_package, _package.Name, _package.Name, _package.Version)
		_package.Metadata = NpmMetadata

		Packages = append(Packages, _package)

	}
	return nil
}

// Parse lock content
func readNpmLockContent(location *model.Location) error {

	file, err := os.ReadFile(location.Path)
	if err != nil {
		return err
	}

	// Validate if file is valid JSON
	if !json.Valid(file) {
		return nil
	}

	if err = json.Unmarshal(file, &NPMLockMetadata); err != nil {
		return err
	}

	if len(NPMLockMetadata.Dependencies) > 0 {
		for name, cPackage := range NPMLockMetadata.Dependencies {
			_package := new(model.Package)
			_package.ID = uuid.NewString()
			_package.Locations = append(_package.Locations, model.Location{
				Path:      TrimUntilLayer(*location),
				LayerHash: location.LayerHash,
			})

			// // init npm data
			_package.Name = name
			_package.Version = cPackage.Version
			_package.Type = npm
			_package.Path = name

			// //parseURL
			parseNpmPackageURL(_package)
			cpe.NewCPE23(_package, _package.Name, _package.Name, _package.Version)
			_package.Metadata = cPackage

			Packages = append(Packages, _package)

		}
	}

	return nil
}

// Parse yarn lock content
func readYarnLockContent(location *model.Location) error {

	file, err := os.Open(location.Path)
	if err != nil {
		return err
	}
	metadata := make(NpmLockMetadata)
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
			_package := new(model.Package)
			_package.ID = uuid.NewString()
			_package.Type = npm
			_package.Name = metadata["Name"].(string)
			_package.Path = metadata["Name"].(string)

			if metadata["Version"] != nil {
				_package.Version = metadata["Version"].(string)
			}

			_package.Locations = append(_package.Locations, model.Location{
				Path:      TrimUntilLayer(*location),
				LayerHash: location.LayerHash,
			})
			parseNpmPackageURL(_package)
			cpe.NewCPE23(_package, _package.Name, _package.Name, _package.Version)
			_package.Metadata = metadata
			Packages = append(Packages, _package)
			metadata = NpmLockMetadata{}
		}
	}

	return nil
}

// Parse PURL
func parseNpmPackageURL(_package *model.Package) {
	_package.PURL = model.PURL(scheme + ":" + npm + "/" + _package.Name + "@" + _package.Version)
}
