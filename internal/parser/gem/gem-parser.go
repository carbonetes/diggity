package gem

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/file"
	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/parser/bom"
	"github.com/carbonetes/diggity/internal/parser/util"

	"github.com/google/uuid"
)

const (
	gemPackage = ".gemspec"
	gem        = "gem"
	spec       = "specifications"
	lockFile   = "Gemfile.lock"
)

// Metadata  metadata
type Metadata map[string]interface{}

// FindGemPackagesFromContent Find gem packages in the file contents
func FindGemPackagesFromContent() {
	if util.ParserEnabled(gem) {
		for _, content := range file.Contents {
			if strings.Contains(content.Path, gemPackage) && strings.Contains(content.Path, spec) {
				if err := readGemContent(content); err != nil {
					err = errors.New("gem-parser: " + err.Error())
					bom.Errors = append(bom.Errors, &err)
				}
			} else if strings.Contains(content.Path, lockFile) {
				if err := readGemLockContent(content); err != nil {
					err = errors.New("gem-parser: " + err.Error())
					bom.Errors = append(bom.Errors, &err)
				}
			}
		}
	}
	defer bom.WG.Done()
}

// Parse gem lock content
func readGemLockContent(location *model.Location) error {
	gemFile, err := os.Open(location.Path)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(gemFile)
	for scanner.Scan() {
		keyValue := scanner.Text()
		trimedKeyValue := strings.TrimSpace(keyValue)

		if len(keyValue) > 1 && keyValue[0] != ' ' {
			continue
		}

		if isKeyValueValid(keyValue) {
			stringArray := strings.Fields(trimedKeyValue)
			if len(stringArray) == 2 {
				_package := new(model.Package)
				_package.ID = uuid.NewString()
				_package.Name = stringArray[0]
				_package.Type = gem
				_package.Path = stringArray[0]
				_package.Version = strings.Trim(stringArray[1], "()")
				//generate cpe
				cpe.NewCPE23(_package, _package.Name, _package.Name, _package.Version)
				//generate and trim path
				_package.Locations = append(_package.Locations, model.Location{
					Path:      util.TrimUntilLayer(*location),
					LayerHash: location.LayerHash,
				})

				bom.Packages = append(bom.Packages, _package)
			}
		}
	}
	return nil
}

// Check if key value is valid
func isKeyValueValid(keyValue string) bool {
	if len(keyValue) < 5 {
		return false
	}
	return strings.Count(keyValue[:5], " ") == 4
}

// Read file contents
func readGemContent(location *model.Location) error {
	gemFile, err := os.Open(location.Path)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(gemFile)

	var value string
	var attribute string
	var previousAttribute string

	metadata := make(Metadata)

	for scanner.Scan() {
		keyValue := scanner.Text()

		if strings.Contains(keyValue, "=") {
			keyValues := strings.SplitN(keyValue, "=", 2)
			attribute = keyValues[0]
			value = keyValues[1]

			//check if attribute is invalid - set to null if invalid
			if strings.Contains(attribute, "%") || strings.Contains(attribute, "if Gem") {
				//clear attribute
				attribute = ""
			}
		} else {
			value = strings.TrimSpace(value + keyValue)
			attribute = previousAttribute
		}

		if len(attribute) > 0 && attribute != " " {
			attribute = strings.ReplaceAll(attribute, " ", "")
			attribute = strings.Replace(attribute, "s.", "", -1)
			value = strings.Replace(value, "\r\n", "", -1)
			value = strings.ReplaceAll(value, ".freeze", "")
			metadata[attribute] = strings.ReplaceAll(value, "\"", "")
			metadata[attribute] = strings.TrimSpace(metadata[attribute].(string))
		}

		previousAttribute = attribute
	}
	if len(metadata) > 0 {
		_package := new(model.Package)
		_package.ID = uuid.NewString()
		_package.Type = gem

		//generate and trim path
		_package.Locations = append(_package.Locations, model.Location{
			Path:      util.TrimUntilLayer(*location),
			LayerHash: location.LayerHash,
		})

		initGemPackages(_package, metadata)
		bom.Packages = append(bom.Packages, _package)
	}

	return nil
}

// Initialize package
func initGemPackages(p *model.Package, metadata Metadata) *model.Package {

	re := regexp.MustCompile(`[^\w^,^ ]`)

	var licenses = make([]string, 0)
	// var authors []string = make([]string, 0)
	p.Name = metadata["name"].(string)
	p.Path = metadata["name"].(string)
	p.Version = metadata["version"].(string)
	if val, ok := metadata["description"].(string); ok {
		p.Description = val
	}
	if val, ok := metadata["licenses"].(string); ok {
		tmpLicenses := re.ReplaceAllString(val, "")
		licenses = append(licenses, tmpLicenses)
	}
	p.Licenses = licenses
	p.Type = gem

	//parseURL
	parseGemPackageURL(p)

	//check if metadata key is exist. if exist delete key to avoid duplicates
	if _, ok := metadata["metadata"].(string); ok {
		delete(metadata, "metadata")
	}

	//check if authors exists
	if val, ok := metadata["authors"].(string); ok {
		tmpAuthors := re.ReplaceAllString(val, "")
		if strings.Contains(tmpAuthors, ",") {
			arrAuthors := strings.Split(tmpAuthors, ", ")
			metadata["authors"] = arrAuthors
			for _, tmpAuthor := range arrAuthors {
				cpe.NewCPE23(p, strings.TrimSpace(tmpAuthor), p.Name, p.Version)
			}

		} else {
			var authors = make([]string, 0)
			authors = append(authors, tmpAuthors)
			metadata["authors"] = authors
			cpe.NewCPE23(p, strings.TrimSpace(tmpAuthors), p.Name, p.Version)
		}
	}

	//check if files exists
	if val, ok := metadata["files"].(string); ok {
		tmpFiles := re.ReplaceAllString(val, "")
		if strings.Contains(tmpFiles, ",") {
			metadata["files"] = strings.Split(tmpFiles, ", ")
		} else {
			var files = make([]string, 0)
			files = append(files, tmpFiles)
			metadata["files"] = files
		}
	}
	metadata["licenses"] = licenses
	p.Metadata = metadata

	return p
}

// Parse PURL
func parseGemPackageURL(_package *model.Package) {
	_package.PURL = model.PURL("pkg" + ":" + gem + "/" + _package.Name + "@" + _package.Version)
}
