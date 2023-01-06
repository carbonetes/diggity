package parser

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/file"
	"github.com/carbonetes/diggity/internal/model"

	"github.com/google/uuid"
)

const (
	pythonPackage = "METADATA"
	pythonRecord  = "RECORD"
	pythonEgg     = ".egg-info"
	pip           = "python"
	unknownField  = "UNKNOWN"
)

// PythonMetadata  metadata
type PythonMetadata map[string]interface{}

// FindPythonPackagesFromContent - Find python packages in the file contents
func FindPythonPackagesFromContent() {
	if parserEnabled(pip) {
		for _, content := range file.Contents {
			if strings.Contains(content.Path, pythonPackage) || strings.Contains(content.Path, pythonEgg) {
				if err := readPythonContent(content); err != nil {
					err = errors.New("python-parser: " + err.Error())
					Errors = append(Errors, &err)
				}
			}
		}
	}
	defer WG.Done()
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
	// reg := regexp.MustCompile(`[^\w^,^ ]`)

	metadata := make(PythonMetadata)

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
		_package := new(model.Package)
		_package.ID = uuid.NewString()
		_package.Type = pip
		_package.Locations = append(_package.Locations, model.Location{
			Path:      TrimUntilLayer(*location),
			LayerHash: location.LayerHash,
		})
		initPythonPackages(_package, metadata, location)
		Packages = append(Packages, _package)
	}

	return nil
}

// Initialize python package
func initPythonPackages(p *model.Package, metadata map[string]interface{}, location *model.Location) *model.Package {

	p.Name = metadata["Name"].(string)
	p.Version = metadata["Version"].(string)
	p.Path = metadata["Name"].(string)

	//check first if description exist in metadata
	if val, ok := metadata["description"].(string); ok {
		p.Description = val
	}

	//check first if license exist in metadata
	if val, ok := metadata["License"]; ok {
		p.Licenses = append(p.Licenses, val.(string))
	}

	p.Type = pip

	//parseURL
	parsePythonPackageURL(p)
	filesPath := strings.Split(location.Path, pythonPackage)[0]
	filesPath = filesPath + pythonRecord
	err := parseMetadataFiles(metadata, filesPath)
	if _, ok := metadata["Files"]; ok && err == nil {
		tmpLocation := new(model.Location)
		tmpLocation.LayerHash = location.LayerHash
		tmpLocation.Path = filesPath
		p.Locations = append(p.Locations, model.Location{
			Path:      TrimUntilLayer(*tmpLocation),
			LayerHash: location.LayerHash,
		})
	}
	p.Metadata = metadata

	//parse CPE
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
func parseMetadataFiles(m PythonMetadata, path string) error {
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
				pathValue = ""
				algorithm = ""
				algorithmValue = ""
				valueSize = ""
			}

		}
	}
	m["Files"] = files
	return nil
}

// Parse PURL
func parsePythonPackageURL(_package *model.Package) {
	_package.PURL = model.PURL(scheme + ":" + pip + "/" + _package.Name + "@" + _package.Version)
}
