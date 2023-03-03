package python

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
	pythonPackage = "METADATA"
	pythonRecord  = "RECORD"
	pythonEgg     = ".egg-info"
	pip           = "python"
	pypi          = "pypi"
	unknownField  = "UNKNOWN"
	poetry        = "poetry.lock"
	poetryPackage = "[[package]]"
	requirements  = "requirements"
	txt           = ".txt"
	poetryFiles   = "files = ["
	fileHashKey   = "sha256"
	packageTag    = "[[package]]"
)

// Metadata  metadata
type Metadata map[string]interface{}

// FindPythonPackagesFromContent - Find python packages in the file contents
func FindPythonPackagesFromContent() {
	if util.ParserEnabled(pip) {
		for _, content := range file.Contents {
			if strings.Contains(content.Path, pythonPackage) || strings.Contains(content.Path, pythonEgg) {
				if err := readPythonContent(content); err != nil {
					err = errors.New("python-parser: " + err.Error())
					bom.Errors = append(bom.Errors, &err)
				}
			}
			if filepath.Base(content.Path) == poetry {
				if err := readPoetryContent(content); err != nil {
					err = errors.New("python-parser: " + err.Error())
					bom.Errors = append(bom.Errors, &err)
				}
			}
			if strings.Contains(filepath.Base(content.Path), requirements) &&
				strings.Contains(filepath.Base(content.Path), txt) {
				if err := readRequirementsContent(content); err != nil {
					err = errors.New("python-parser: " + err.Error())
					bom.Errors = append(bom.Errors, &err)
				}
			}
		}
	}
	defer bom.WG.Done()
}

// Read file contents
func readPythonContent(location *model.Location) error {
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

		if strings.Contains(keyValue, ": ") {
			keyValues := strings.SplitN(keyValue, ": ", 2)
			attribute = keyValues[0]
			if strings.Contains(attribute, " ") {
				//clear attribute
				attribute = ""
			}
			value = keyValues[1]
		} else {
			value = strings.TrimSpace(value + keyValue)
			attribute = previousAttribute
		}

		if len(attribute) > 0 && attribute != " " {
			metadata[attribute] = strings.Replace(value, "\r\n", "", -1)
			metadata[attribute] = strings.Replace(value, "\r ", "", -1)
		}

		previousAttribute = attribute
	}
	if len(metadata) > 0 && metadata["Name"] != nil {
		bom.Packages = append(bom.Packages, initPythonPackages(metadata, location))
	}

	return nil
}

// Read poetry.lock contents
func readPoetryContent(location *model.Location) error {
	// Read poetry.lock file
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

	isFile := false // check for file metadata
	fileHash := []map[string]string{}

	// Iterate through key value pairs
	for scanner.Scan() {
		keyValue := scanner.Text()

		// Parse files metadata, if any
		if !*bom.Arguments.DisableFileListing {
			if strings.Contains(keyValue, poetryFiles) {
				isFile = true
				continue
			}

			if isFile {
				// assign file metadata and reset
				if strings.TrimSpace(keyValue) == "]" {
					metadata["files"] = fileHash
					fileHash = []map[string]string{}
					isFile = false
					continue
				}
				// skip invalid key values
				if !strings.Contains(keyValue, "file") && !strings.Contains(keyValue, "hash") {
					continue
				}

				fileHash = append(fileHash, poetryFileMetadata(keyValue))
				continue
			}
		} else if strings.Contains(keyValue, "file") {
			continue
		}

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
			// cleanup python-versions metadata
			if _, ok := metadata["python-versions"]; ok {
				metadata["python-versions"] = strings.Replace(metadata["python-versions"].(string), "]", "", -1)
			}
			// init poetry data
			if metadata["name"] != nil {
				bom.Packages = append(bom.Packages, initPythonPackages(metadata, location))
			}

			// Reset metadata
			metadata = make(Metadata)
		}
	}

	// Parse packages before EOF
	if metadata["name"] != nil {
		bom.Packages = append(bom.Packages, initPythonPackages(metadata, location))
	}

	return nil
}

// Read requirements.txt contents
func readRequirementsContent(location *model.Location) error {
	// Read requirements.txt file
	file, err := os.Open(location.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Iterate through requirements
	for scanner.Scan() {
		req := scanner.Text()

		// Remove Comments
		if strings.Contains(req, "#") {
			req = strings.Split(req, "#")[0]
		}

		if strings.Contains(req, "==") {
			name, version := parseRequirements(req)

			if name != "" && !strings.Contains(name, ";") {
				metadata := make(Metadata)
				metadata["name"] = name
				metadata["version"] = version
				bom.Packages = append(bom.Packages, initPythonPackages(metadata, location))
			}
		}

	}

	return nil
}

// Initialize python package
func initPythonPackages(metadata map[string]interface{}, location *model.Location) *model.Package {
	p := new(model.Package)
	p.ID = uuid.NewString()
	p.Type = pip
	p.Locations = append(p.Locations, model.Location{
		Path:      util.TrimUntilLayer(*location),
		LayerHash: location.LayerHash,
	})

	// parse name and version based on metadata
	if _, ok := metadata["Name"]; ok {
		p.Name = metadata["Name"].(string)
		p.Version = metadata["Version"].(string)
		p.Path = metadata["Name"].(string)
	} else {
		p.Name = metadata["name"].(string)
		p.Version = metadata["version"].(string)
		p.Path = metadata["name"].(string)
	}

	// check first if description exist in metadata
	if val, ok := metadata["description"].(string); ok {
		p.Description = val
	}

	// check first if license exist in metadata
	if val, ok := metadata["License"]; ok {
		p.Licenses = append(p.Licenses, val.(string))
	} else {
		p.Licenses = []string{}
	}

	p.Type = pip

	// parse PURL
	parsePythonPackageURL(p)
	filesPath := strings.Split(location.Path, pythonPackage)[0]
	filesPath = filesPath + pythonRecord
	err := parseMetadataFiles(metadata, filesPath)
	if _, ok := metadata["Files"]; ok && err == nil {
		tmpLocation := new(model.Location)
		tmpLocation.LayerHash = location.LayerHash
		tmpLocation.Path = filesPath
		p.Locations = append(p.Locations, model.Location{
			Path:      util.TrimUntilLayer(*tmpLocation),
			LayerHash: location.LayerHash,
		})
	}
	p.Metadata = metadata

	// parse CPE
	if val, ok := metadata["Author"].(string); ok {
		if val == unknownField {
			val = p.Name
		}
		cpe.NewCPE23(p, strings.TrimSpace(val), p.Name, p.Version)
	} else {
		cpe.NewCPE23(p, p.Name, p.Name, p.Version)
	}

	return p
}

// Parse python metadata
func parseMetadataFiles(m Metadata, path string) error {
	var mapValue = map[string]interface{}{}
	var files []map[string]interface{}
	var finalValue = map[string]interface{}{}
	var pathValue string
	var algorithm string
	var algorithmValue string
	var valueSize string
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("file does not exists: %s", path)
	}

	fileinfo, _ := os.ReadFile(path)

	lines := strings.Split(string(fileinfo), "\n")
	for _, line := range lines {

		if strings.Contains(line, ",") {
			keyValues := strings.Split(line, ",")
			pathValue = keyValues[0]
			tmpValue := keyValues[1]
			if tmpValue != "" {
				tmpSplitValue := strings.Split(tmpValue, "=")
				algorithm = tmpSplitValue[0]
				algorithmValue = tmpSplitValue[1]
			}
			valueSize = keyValues[2]

			if pathValue != "" {
				finalValue["path"] = pathValue
				mapValue["value"] = algorithmValue
				mapValue["algorithm"] = algorithm

				//ignore digest if algorithm is blank
				if algorithm != "" {
					finalValue["digest"] = mapValue
				}
				//ignore valueSize if blank
				if valueSize != "" {
					valueSize = strings.Replace(valueSize, "\r", "", -1)
					finalValue["size"] = valueSize
				}
				files = append(files, finalValue)
				algorithm = ""
				algorithmValue = ""
			}

		}
	}
	m["Files"] = files
	return nil
}

// Parse PURL
func parsePythonPackageURL(_package *model.Package) {
	_package.PURL = model.PURL("pkg" + ":" + pypi + "/" + _package.Name + "@" + _package.Version)
}

// Parse requirements metadata
func parseRequirements(req string) (name string, version string) {
	reqMetadata := strings.Split(req, "==")
	versionMetadata := strings.TrimSpace(reqMetadata[1])

	name = strings.TrimSpace(reqMetadata[0])
	version = strings.Split(versionMetadata, " ")[0]

	return name, version
}

// Parse poetry file
func poetryFileMetadata(file string) map[string]string {
	fileHash := make(map[string]string)
	r := regexp.MustCompile(`"(.*?)"`)

	for _, fh := range r.FindAllString(file, -1) {
		// assign to hash if contains sha256
		if strings.Contains(fh, fileHashKey) {
			fileHash["hash"] = strings.Replace(fh, `"`, "", -1)
		} else {
			fileHash["file"] = strings.Replace(fh, `"`, "", -1)
		}
	}

	return fileHash
}
