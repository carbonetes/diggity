package python

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

const (
	pythonRecord = "RECORD"

	pypi           = "pypi"
	unknownField   = "UNKNOWN"
	pythonVersions = "python_versions"
	poetryPackage  = "[[package]]"
	txt            = ".txt"
	poetryFiles    = "files = ["
	fileHashKey    = "sha256"
	packageTag     = "[[package]]"
)

const parserErr = "python-parser: "

// Read file contents
func readPythonContent(location *model.Location, req *bom.ParserRequirements) {
	file, err := os.Open(location.Path)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	defer file.Close()

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
		*req.SBOM.Packages = append(*req.SBOM.Packages, *initPythonPackages(metadata, location))
	}
}

//TODO: rework parser to reduce code complexity

// Read poetry.lock contents
func readPoetryContent(location *model.Location, req *bom.ParserRequirements) {
	// Read poetry.lock file
	file, err := os.Open(location.Path)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}
	defer file.Close()

	metadata := make(Metadata)
	scanner := bufio.NewScanner(file)

	var value, attribute, previousAttribute string

	isFile := false // check for file metadata
	fileHash := []map[string]string{}

	// Iterate through key value pairs
	for scanner.Scan() {
		keyValue := scanner.Text()

		// Parse files metadata, if any
		if !*req.Arguments.DisableFileListing {
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
			if _, ok := metadata[pythonVersions]; ok {
				metadata[pythonVersions] = strings.Replace(metadata[pythonVersions].(string), "]", "", -1)
			}
			// init poetry data
			if metadata["name"] != nil {
				*req.SBOM.Packages = append(*req.SBOM.Packages, *initPythonPackages(metadata, location))
			}

			// Reset metadata
			metadata = make(Metadata)
		}
	}

	// Parse packages before EOF
	if metadata["name"] != nil {
		*req.SBOM.Packages = append(*req.SBOM.Packages, *initPythonPackages(metadata, location))
	}
}

// Read requirements.txt contents
func readRequirementsContent(location *model.Location, req *bom.ParserRequirements) {
	// Read requirements.txt file
	file, err := os.Open(location.Path)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Iterate through requirements
	for scanner.Scan() {
		text := scanner.Text()

		// Remove Comments
		if strings.Contains(text, "#") {
			text = strings.Split(text, "#")[0]
		}

		if !strings.Contains(text, "==") {
			continue
		}

		name, version := parseRequirements(text)

		if name == "" && strings.Contains(name, ";") {
			continue
		}
		metadata := make(Metadata)
		metadata["name"] = name
		metadata["version"] = version

		*req.SBOM.Packages = append(*req.SBOM.Packages, *initPythonPackages(metadata, location))

	}
}
