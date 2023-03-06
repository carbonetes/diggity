package cargo

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"regexp"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/file"
	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/model/metadata"
	"github.com/carbonetes/diggity/internal/parser/bom"
	"github.com/carbonetes/diggity/internal/parser/util"

	"strings"

	"github.com/google/uuid"
)

const (
	rust       = "rust"
	rustCrate  = "rust-crate"
	cargo      = "cargo"
	cargoLock  = "Cargo.lock"
	packageTag = "[[package]]"
)

// Metadata cargo metadata
type Metadata map[string]interface{}

// FindCargoPackagesFromContent checks for cargo.lock files in the file contents
func FindCargoPackagesFromContent() {
	if util.ParserEnabled(rust) {
		for _, content := range file.Contents {
			if filepath.Base(content.Path) == cargoLock {
				if err := readCargoContent(content); err != nil {
					err = errors.New("cargo-parser: " + err.Error())
					bom.Errors = append(bom.Errors, &err)
				}
			}
		}
	}
	defer bom.WG.Done()
}

// Read Cargo.lock package information
func readCargoContent(location *model.Location) error {
	// Read Cargo.lock file
	file, err := os.Open(location.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	metadata := make(Metadata)
	scanner := bufio.NewScanner(file)

	var value string
	var attribute string
	var previousAttribute string

	// Iterate through key value pairs
	for scanner.Scan() {
		keyValue := scanner.Text()

		if strings.Contains(keyValue, "=") {
			keyValues := strings.SplitN(keyValue, "=", 2)
			attribute = util.FormatLockKeyVal(keyValues[0])
			value = util.FormatLockKeyVal(keyValues[1])

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

		// Packages delimited by line breaks or [[package]] tag
		if len(keyValue) <= 1 || keyValue == packageTag {
			// init cargo data
			if metadata["name"] != nil {
				bom.Packages = append(bom.Packages, initRustPackage(location, metadata))
			}

			// Reset metadata
			metadata = make(Metadata)
		}
	}

	// Parse packages before EOF
	if metadata["name"] != nil {
		bom.Packages = append(bom.Packages, initRustPackage(location, metadata))
	}

	return nil
}

// Init Cargo Package
func initRustPackage(location *model.Location, metadata Metadata) *model.Package {
	_package := new(model.Package)
	_package.ID = uuid.NewString()
	_package.Name = metadata["name"].(string)
	_package.Version = metadata["version"].(string)
	_package.Path = _package.Name
	_package.Type = rustCrate
	_package.Locations = append(_package.Locations, model.Location{
		Path:      util.TrimUntilLayer(*location),
		LayerHash: location.LayerHash,
	})
	_package.Licenses = []string{}

	// get purl
	parseRustPackageURL(_package)

	// get CPEs
	cpe.NewCPE23(_package, "", _package.Name, _package.Version)

	// fill metadata
	initCargoMetadata(_package, metadata)

	return _package
}

// Init Cargo Metadata
func initCargoMetadata(p *model.Package, m Metadata) {
	source := ""
	checksum := ""
	deps := []string{}

	// Check if metadata exists
	if m["source"] != nil {
		source = m["source"].(string)
	}
	if m["checksum"] != nil {
		checksum = m["checksum"].(string)
	}
	if m["dependencies"] != nil {
		deps = formatDependencies(m["dependencies"].(string))
	}

	p.Metadata = metadata.CargoMetadata{
		Name:         m["name"].(string),
		Version:      m["version"].(string),
		Source:       source,
		Checksum:     checksum,
		Dependencies: deps,
	}
}

// Parse PURL
func parseRustPackageURL(_package *model.Package) {
	_package.PURL = model.PURL("pkg" + ":" + cargo + "/" + _package.Name + "@" + _package.Version)
}

// Format Dependencies Metadata
func formatDependencies(depsString string) (deps []string) {
	r := regexp.MustCompile(`"(.*?)"`)
	for _, d := range r.FindAllString(depsString, -1) {
		deps = append(deps, util.FormatLockKeyVal(d))
	}
	return deps
}
