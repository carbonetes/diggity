package parser

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/docker"
	"github.com/carbonetes/diggity/internal/file"
	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/model/metadata"

	"github.com/google/uuid"
	"golang.org/x/exp/maps"
)

const (
	pomFileName      = "pom.xml"
	manifestFile     = "MANIFEST.MF"
	propertiesFile   = "pom.properties"
	scheme           = "pkg"
	java             = "java"
	jarPackagesRegex = `(?:\.jar)$|(?:\.ear)$|(?:\.war)$|(?:\.jpi)$|(?:\.hpi)$`
)

type (
	// JavaMetadata java metadata
	JavaMetadata map[string]map[string]string
	// JavaManifest java manifest
	JavaManifest map[string]string
)

var (
	// JavaPomXML pom metadata
	JavaPomXML metadata.Project
	// Result java metadata
	Result              = make(map[string]*model.Package, 0)
	nameAndVersionRegex = regexp.MustCompile(`(?Ui)^(?P<name>(?:[[:alpha:]][[:word:].]*(?:\.[[:alpha:]][[:word:].]*)*-?)+)(?:-(?P<version>(?:\d.*|(?:build\d*.*)|(?:rc?\d+(?:^[[:alpha:]].*)?))))?$`)
)

// FindJavaPackagesFromContent checks for jar files in the file contents
func FindJavaPackagesFromContent() {
	if parserEnabled(java) {
		for _, content := range file.Contents {
			if match := regexp.MustCompile(jarPackagesRegex).FindString(content.Path); len(match) > 0 {
				if err := extractJarFile(content); err != nil {
					err = errors.New("java-parser: " + err.Error())
					Errors = append(Errors, &err)
				}
			} else if strings.Contains(content.Path, pomFileName) {
				if err := parsePomXML(*content, content.Path); err != nil {
					err = errors.New("java-parser: " + err.Error())
					Errors = append(Errors, &err)
				}
			}
		}
		Packages = append(Packages, maps.Values(Result)...)
	}
	defer WG.Done()
}

// Extract jar files
func extractJarFile(location *model.Location) error {

	buff, err := os.ReadFile(location.Path)
	if err != nil {
		return err
	}
	reader, err := zip.NewReader(bytes.NewReader(buff), int64(len(buff)))
	if err != nil {
		return err
	}
	var dependencies []*zip.File
	for _, zipFile := range reader.File {
		if zipFile == nil {
			break
		}

		if match := regexp.MustCompile(jarPackagesRegex).FindString(zipFile.Name); len(match) > 0 {
			dependencies = append(dependencies, zipFile)
		} else if strings.Contains(zipFile.Name, pomFileName) {
			file, err := os.Create(filepath.Join(docker.Dir(), pomFileName))
			if err != nil {
				panic(err)
			}
			reader, err := zipFile.Open()
			if err != nil {
				panic(err)
			}
			io.Copy(file, reader)
			parsePomXML(model.Location{
				Path:      file.Name(),
				LayerHash: location.LayerHash,
			}, location.Path)
			if err := file.Close(); err != nil {
				return err
			}
			defer os.Remove(file.Name())
		}
		var metadataFile *zip.File = nil
		var pomPropertiesFile *zip.File = nil

		metadataFile, pomPropertiesFile = checkZipfile(zipFile, metadataFile, pomPropertiesFile)

		if metadataFile != nil || pomPropertiesFile != nil {
			paths := strings.Split(TrimUntilLayer(*location), string(os.PathSeparator))
			initPackage(paths[len(paths)-1], location, metadataFile, pomPropertiesFile)
			pomPropertiesFile = nil
			metadataFile = nil
		}
	}

	err = parseJarFiles(dependencies, location)
	return err
}

// Init java package
func initPackage(name string, location *model.Location, manifestFile *zip.File, pomPropertiesFile *zip.File) error {
	endOfFile := regexp.MustCompile(jarPackagesRegex)
	_package := new(model.Package)
	_package.Metadata = JavaMetadata{}
	_package.ID = uuid.NewString()
	_package.Path = name
	_package.Type = java
	paths := strings.Split(location.Path, string(os.PathSeparator))
	_package.Locations = append(_package.Locations, model.Location{
		Path:      paths[len(paths)-1],
		LayerHash: location.LayerHash,
	})
	splitName := strings.Split(name, "/")
	fileName := splitName[len(splitName)-1]
	if manifestFile != nil {
		parseJavaManifest(manifestFile, _package)
		parseLicenses(_package)
	}

	vendor := ""
	product := ""
	version := ""

	if pomPropertiesFile != nil {
		reader, err := pomPropertiesFile.Open()
		if err != nil {
			return err
		}
		data, err := io.ReadAll(reader)

		if err != nil {
			return err
		}
		parsePomProperties(string(data), _package, pomPropertiesFile.Name)
		if err := reader.Close(); err != nil {
			return err
		}
	}
	if pomPropertiesFile != nil {
		vendor = _package.Metadata.(JavaMetadata)["PomProperties"]["artifactId"]
		version = _package.Metadata.(JavaMetadata)["PomProperties"]["version"]
		product = vendor
	} else {
		version = formatVersionMetadata(_package, _package.Metadata.(JavaMetadata)["Manifest"]["Implementation-Version"])
		regex, err := regexp.Compile(`-\s*(\d+)`)
		if err != nil {
			return err
		}
		names := regex.Split(fileName, -1)
		vendor = names[0]
		vendor = endOfFile.ReplaceAllString(vendor, "")
		product = vendor
	}

	if len(version) == 0 {
		matches := nameAndVersionRegex.FindStringSubmatch(strings.Replace(fileName, ".jar", "", 1))
		if len(matches) >= 2 {
			vendor = matches[1]
			version = matches[2]
		} else {
			version = strings.Replace(
				strings.Replace(
					strings.Replace(fileName, vendor, "", 1), "-", "", 1), ".jar", "", 1)
		}
		product = vendor
	}

	_package.Name = vendor
	_package.Version = version
	_package.Description = _package.Metadata.(JavaMetadata)["Manifest"]["Bundle-Description"]

	parseJavaURL(_package)

	cpe.NewCPE23(_package, vendor, product, version)
	// additional CPEs if
	if _package.Metadata.(JavaMetadata)["PomProperties"] != nil {
		vendor = _package.Metadata.(JavaMetadata)["PomProperties"]["groupId"]
		if len(vendor) > 0 {
			generateAdditionalCPE(vendor, product, version, _package)
		}
	}
	if _package.Metadata.(JavaMetadata)["Manifest"] != nil {
		vendor = _package.Metadata.(JavaMetadata)["Manifest"]["Automatic-Module-Name"]
		if len(vendor) > 0 {
			cpe.NewCPE23(_package, vendor, product, version)
		}
		vendor = _package.Metadata.(JavaMetadata)["Manifest"]["Bundle-SymbolicName"]
		generateAdditionalCPE(vendor, product, version, _package)
	}

	if sourceIsDir() || _package.Name != "" && _package.Version != "" {
		checkPackage(_package, location.LayerHash)
	}
	return nil
}

func checkPackage(_package *model.Package, layerHash string) {
	if _, exists := Result[_package.Name+":"+_package.Version+":"+layerHash]; !exists {
		Result[_package.Name+":"+_package.Version+":"+layerHash] = _package
	} else {
		_tmpPackage := Result[_package.Name+":"+_package.Version+":"+layerHash]
		if _package.Metadata.(JavaMetadata)["Manifest"] != nil {
			_tmpPackage.Metadata.(JavaMetadata)["Manifest"] = _package.Metadata.(JavaMetadata)["Manifest"]
			_tmpPackage.CPEs = append(_tmpPackage.CPEs, _package.CPEs...)
		}
		if _package.Metadata.(JavaMetadata)["PomProperties"] != nil {
			_tmpPackage.Metadata.(JavaMetadata)["PomProperties"] = _package.Metadata.(JavaMetadata)["PomProperties"]
			_tmpPackage.CPEs = append(_tmpPackage.CPEs, _package.CPEs...)
		}
		if _package.Metadata.(JavaMetadata)["PomProject"] != nil {
			_tmpPackage.Metadata.(JavaMetadata)["PomProject"] = _package.Metadata.(JavaMetadata)["PomProject"]
			_tmpPackage.CPEs = append(_tmpPackage.CPEs, _package.CPEs...)

		}
		_tmpPackage.CPEs = cpe.RemoveDuplicateCPES(_tmpPackage.CPEs)
		Result[_package.Name+":"+_package.Version+":"+layerHash] = _tmpPackage
	}

}

// Parse jar files
func parseJarFiles(dependencies []*zip.File, location *model.Location) error {

	if err := os.Mkdir(docker.Dir(), fs.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}

	for _, dependency := range dependencies {

		splitName := strings.Split(dependency.Name, "/")
		fileName := splitName[len(splitName)-1]
		jarFile, err := os.Create(filepath.Join(docker.Dir(), fileName))

		if err != nil {
			return err
		}

		reader, err := dependency.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(jarFile, reader); err != nil {
			return err
		}

		err = findManifestAndPomPropertiesFromDependencyJarFile(jarFile, location, dependency.Name)
		if err != nil {
			return err
		}

		if err := jarFile.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Generate additional CPEs for java packages
func generateAdditionalCPE(vendor string, product string, version string, _package *model.Package) {
	if len(vendor) > 0 {
		if strings.Contains(vendor, ".") {
			for _, v := range strings.Split(vendor, ".") {
				tldsRegex := `(com|org|io|edu|net|edu|gov|mil|the\ |a\ |an\ )(?:\b|')`
				if !regexp.MustCompile(tldsRegex).MatchString(v) {
					cpe.NewCPE23(_package, v, product, version)
				}
			}
		} else {
			cpe.NewCPE23(_package, vendor, product, version)
		}
	}
}

// Find the manifest and pom properties files from the dependency jar file
func findManifestAndPomPropertiesFromDependencyJarFile(jarFile *os.File, location *model.Location, name string) error {

	buff, err := os.ReadFile(jarFile.Name())
	if err != nil {
		return err
	}

	reader, err := zip.NewReader(bytes.NewReader(buff), int64(len(buff)))
	if err != nil {
		return err
	}

	var dependencies []*zip.File
	var metadataFile *zip.File = nil
	var pomPropertiesFile *zip.File = nil

	for _, zipFile := range reader.File {

		if match := regexp.MustCompile(jarPackagesRegex).FindString(zipFile.Name); len(match) > 0 {
			dependencies = append(dependencies, zipFile)
		} else if strings.Contains(zipFile.Name, pomFileName) {
			file, err := os.Create(filepath.Join(docker.Dir(), pomFileName))
			if err != nil {
				panic(err)
			}
			reader, err := zipFile.Open()
			if err != nil {
				panic(err)
			}
			io.Copy(file, reader)
			paths := strings.Split(jarFile.Name(), "/")
			parsePomXML(model.Location{
				Path:      file.Name(),
				LayerHash: location.LayerHash,
			}, paths[len(paths)-1])
			if err := file.Close(); err != nil {
				return err
			}
			defer os.Remove(file.Name())
		}

		metadataFile, pomPropertiesFile = checkZipfile(zipFile, metadataFile, pomPropertiesFile)

		if metadataFile != nil || pomPropertiesFile != nil {
			initPackage(name, location, metadataFile, pomPropertiesFile)
			pomPropertiesFile = nil
			metadataFile = nil
		}
	}
	parseJarFiles(dependencies, location)
	return nil
}

// Validate zip file
func checkZipfile(zipFile *zip.File, metadataFile *zip.File, pomPropertiesFile *zip.File) (*zip.File, *zip.File) {

	if strings.Contains(zipFile.Name, manifestFile) && metadataFile == nil {
		metadataFile = zipFile
	}
	if strings.Contains(zipFile.Name, propertiesFile) && pomPropertiesFile == nil {
		pomPropertiesFile = zipFile
	}
	return metadataFile, pomPropertiesFile

}

// Parse licenses
func parseLicenses(_package *model.Package) {
	var licenses = make([]string, 0)
	for key, value := range _package.Metadata.(JavaMetadata)["Manifest"] {
		if strings.Contains(key, "Bundle-License") {
			licenses = append(licenses, strings.TrimSpace(value))
		}
	}

	_package.Licenses = licenses
}

// Parse java manifest file
func parseJavaManifest(manifestFile *zip.File, _package *model.Package) error {

	createdManifest, err := os.Create(filepath.Join(docker.Dir(), strings.Replace(manifestFile.Name, "/", "_", -1)))
	if err != nil {
		return err
	}
	defer createdManifest.Close()

	var value string
	var attribute string
	var previousAttribute string

	manifest := make(JavaManifest)
	reader, err := manifestFile.Open()

	if err != nil {
		return err
	}

	defer reader.Close()

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		keyValue := strings.TrimSpace(scanner.Text())
		if strings.Contains(keyValue, ":") && !strings.Contains(keyValue, ":=") {
			keyValues := strings.SplitN(keyValue, ":", 2)
			attribute = keyValues[0]
			value = keyValues[1]
		} else {
			value = strings.TrimSpace(value + keyValue)
			attribute = previousAttribute
		}
		if len(attribute) > 0 && attribute != " " {
			manifest[attribute] = strings.Replace(value, "\r\n", "", -1)
			manifest[attribute] = strings.Replace(value, "\r ", "", -1)
			manifest[attribute] = strings.TrimSpace(manifest[attribute])
		}
		previousAttribute = attribute
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	_package.Metadata.(JavaMetadata)["ManifestLocation"] = JavaManifest{"path": filepath.Join(_package.Path, manifestFile.Name)}
	_package.Metadata.(JavaMetadata)["Manifest"] = manifest
	return nil
}

// Parse pom properties
func parsePomProperties(data string, _package *model.Package, path string) {

	var value string
	var attribute string

	pomProperties := make(JavaManifest)
	pomProperties["location"] = filepath.Join(_package.Name, path)
	pomProperties["name"] = ""

	lines := strings.Split(data, "\n")
	for _, keyValue := range lines {
		if strings.Contains(keyValue, "=") {
			keyValues := strings.Split(keyValue, "=")
			attribute = keyValues[0]
			value = keyValues[1]
		}

		if len(attribute) > 0 && attribute != " " {
			pomProperties[attribute] = strings.Replace(value, "\r\n", "", -1)
			pomProperties[attribute] = strings.Replace(value, "\r ", "", -1)
			pomProperties[attribute] = strings.TrimSpace(pomProperties[attribute])
		}
	}

	_package.Metadata.(JavaMetadata)["PomProperties"] = pomProperties
}

// Parse PURL
func parseJavaURL(_package *model.Package) {
	if _package.Metadata.(JavaMetadata)["PomProperties"] != nil {
		_package.PURL = model.PURL(scheme + ":" + "maven" + "/" + _package.Metadata.(JavaMetadata)["PomProperties"]["groupId"] + "/" + _package.Metadata.(JavaMetadata)["PomProperties"]["artifactId"] + "@" + _package.Version)
	} else {
		_package.PURL = model.PURL(scheme + ":" + "maven" + "/" + _package.Name + "/" + _package.Name + "@" + _package.Version)
	}
}

func parsePomXML(location model.Location, layerPath string) error {
	file, err := os.ReadFile(location.Path)
	if err != nil {
		return err
	}
	if err = xml.Unmarshal(file, &JavaPomXML); err != nil {
		return err
	}

	if len(JavaPomXML.Dependencies) > 0 {
		if sourceIsDir() {
			for _, dep := range JavaPomXML.Dependencies {
				if dep.ArtifactID != "" && !strings.Contains(dep.Version, "$") {
					_package := new(model.Package)
					_package.Metadata = JavaMetadata{}
					_package.ID = uuid.NewString()
					_package.Name = dep.ArtifactID
					_package.Path = TrimUntilLayer(model.Location{
						Path:      layerPath,
						LayerHash: location.LayerHash,
					})
					_package.Version = dep.Version
					_package.Type = java
					paths := strings.Split(location.Path, string(os.PathSeparator))
					_package.Locations = append(_package.Locations, model.Location{
						Path:      paths[len(paths)-1],
						LayerHash: location.LayerHash,
					})
					_package.Metadata.(JavaMetadata)["ManifestLocation"] = JavaManifest{"path": _package.Path}
					_package.Metadata.(JavaMetadata)["PomProject"] = JavaManifest{
						"name":    _package.Name,
						"version": _package.Version,
						"groupID": dep.GroupID,
					}
					parseJavaURL(_package)
					cpe.NewCPE23(_package, _package.Name, _package.Name, _package.Version)
					generateAdditionalCPE(dep.GroupID, _package.Name, _package.Version, _package)
					checkPackage(_package, location.LayerHash)
				}
			}
		} else {
			if _, exists := Result[JavaPomXML.ArtifactID+":"+JavaPomXML.Version+":"+location.LayerHash]; exists {
				_tmpPackage := Result[JavaPomXML.ArtifactID+":"+JavaPomXML.Version+":"+location.LayerHash]
				_tmpPackage.Metadata.(JavaMetadata)["PomProject"] = JavaManifest{
					"name":    JavaPomXML.Name,
					"version": JavaPomXML.Version,
					"groupID": JavaPomXML.GroupID,
				}
				cpe.NewCPE23(_tmpPackage, JavaPomXML.ArtifactID, JavaPomXML.ArtifactID, JavaPomXML.Version)
				generateAdditionalCPE(JavaPomXML.GroupID, JavaPomXML.ArtifactID, JavaPomXML.Version, _tmpPackage)
				Result[JavaPomXML.ArtifactID+":"+JavaPomXML.Version+":"+location.LayerHash] = _tmpPackage
			}
		}
	}
	return nil
}

// Format Versions with other metadata appended
func formatVersionMetadata(p *model.Package, version string) string {
	r := regexp.MustCompile(`[A-Z][a-z]*(\s[A-Z][a-z]*)*`)
	match := r.FindString(version)

	if match != "" && len(match) > 1 {
		versionMetadata := strings.Split(version, match)
		p.Metadata.(JavaMetadata)["Manifest"]["Implementation-Version"] = versionMetadata[0]

		newMetadata := strings.SplitN(match+versionMetadata[1], ":", 2)
		if len(newMetadata) > 1 {
			p.Metadata.(JavaMetadata)["Manifest"][newMetadata[0]] = strings.TrimSpace(newMetadata[1])
		}

		return versionMetadata[0]
	}

	return version
}
