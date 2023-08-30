package npm

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
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

const parserErr string = "npm-parser: "

// Read file contents
func readNpmContent(location *model.Location, req *bom.ParserRequirements) {
	file, err := os.ReadFile(location.Path)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	// Validate if file is valid JSON
	if !json.Valid(file) {
		return
	}

	if err = json.Unmarshal(file, &NpmMetadata); err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	if NpmMetadata.Name == "" {
		return
	}

	pkg := newNpmPackage(&NpmMetadata)
	if pkg == nil || pkg.Name == "" {
		return
	}

	pkg.Locations = append(pkg.Locations, model.Location{
		Path:      util.TrimUntilLayer(*location),
		LayerHash: location.LayerHash,
	})
	*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)

}

// Parse lock content
func readNpmLockContent(location *model.Location, req *bom.ParserRequirements) {

	file, err := os.ReadFile(location.Path)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	// Validate if file is valid JSON
	if !json.Valid(file) {
		return
	}

	if err = json.Unmarshal(file, &NpmLockMetadata); err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	if len(NpmLockMetadata.Dependencies) > 0 {
		for name, metadata := range NpmLockMetadata.Dependencies {
			pkg := newNpmLockPackage(name, &metadata)
			if pkg == nil || pkg.Name == "" {
				continue
			}
			pkg.Locations = append(pkg.Locations, model.Location{
				Path:      util.TrimUntilLayer(*location),
				LayerHash: location.LayerHash,
			})
			*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)

		}
	}
}

// Parse yarn lock content
func readYarnLockContent(location *model.Location, req *bom.ParserRequirements) {
	file, err := readFile(location.Path)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	var value, attribute string
	metadata := make(LockMetadata)
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
			pkg := newYarnLockPackage(&metadata)
			pkg.Locations = append(pkg.Locations, model.Location{
				Path:      util.TrimUntilLayer(*location),
				LayerHash: location.LayerHash,
			})
			*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
		}
	}
}
