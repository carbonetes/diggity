package java

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
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"

	"github.com/google/uuid"
	"golang.org/x/exp/maps"
)

const (
	pomFileName      = "pom.xml"
	manifestFile     = "MANIFEST.MF"
	propertiesFile   = "pom.properties"
	java             = "java"
	jarPackagesRegex = `(?:\.jar)$|(?:\.ear)$|(?:\.war)$|(?:\.jpi)$|(?:\.hpi)$`
)

type (
	// Metadata java metadata
	Metadata map[string]map[string]string
	// Manifest java manifest
	Manifest map[string]string
)

var (
	// JavaPomXML pom metadata
	JavaPomXML metadata.Project
	// Result java metadata
	Result              = make(map[string]model.Package, 0)
	nameAndVersionRegex = regexp.MustCompile(`(?Ui)^(?P<name>(?:[[:alpha:]][[:word:].]*(?:\.[[:alpha:]][[:word:].]*)*-?)+)(?:-(?P<version>(?:\d.*|(?:build\d*.*)|(?:rc?\d+(?:^[[:alpha:]].*)?))))?$`)
)

// FindJavaPackagesFromContent checks for jar files in the file contents
func FindJavaPackagesFromContent(req *bom.ParserRequirements) {
	if util.ParserEnabled(java, req.Arguments.EnabledParsers) {
		for _, content := range *req.Contents {
			if match := regexp.MustCompile(jarPackagesRegex).FindString(content.Path); len(match) > 0 {
				if err := extractJarFile(&content, req.Arguments.Dir, req.DockerTemp); err != nil {
					err = errors.New("java-parser: " + err.Error())
					*req.Errors = append(*req.Errors, err)
				}
			} else if strings.Contains(content.Path, pomFileName) {
				if err := parsePomXML(content, content.Path, req.Arguments.Dir); err != nil {
					err = errors.New("java-parser: " + err.Error())
					*req.Errors = append(*req.Errors, err)
				}
			}
		}
		*req.SBOM.Packages = append(*req.SBOM.Packages, maps.Values(Result)...)
	}
	defer req.WG.Done()
}

// Extract jar files
func extractJarFile(location *model.Location, dir *string, dockerTemp *string) error {

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
		// Skip unsafe files fo extraction
		if strings.Contains(zipFile.Name, "..") {
			continue
		}

		if match := regexp.MustCompile(jarPackagesRegex).FindString(zipFile.Name); len(match) > 0 {
			dependencies = append(dependencies, zipFile)
		} else if strings.Contains(zipFile.Name, pomFileName) {
			file, err := os.Create(filepath.Join(*dockerTemp, pomFileName))
			if err != nil {
				panic(err)
			}
			reader, err := zipFile.Open()
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(file, reader)
			if err != nil {
				return err
			}

			_ = parsePomXML(model.Location{
				Path:      file.Name(),
				LayerHash: location.LayerHash,
			}, location.Path, dir)

			if err := file.Close(); err != nil {
				return err
			}
			defer os.Remove(file.Name())
		}
		var metadataFile *zip.File = nil
		var pomPropertiesFile *zip.File = nil

		metadataFile, pomPropertiesFile = checkZipfile(zipFile, metadataFile, pomPropertiesFile)

		if metadataFile != nil || pomPropertiesFile != nil {
			paths := strings.Split(util.TrimUntilLayer(*location), string(os.PathSeparator))
			_ = initPackage(paths[len(paths)-1], location, metadataFile, pomPropertiesFile, dir, dockerTemp)
		}
	}

	err = parseJarFiles(dependencies, location, dir, dockerTemp)
	return err
}

// Init java package
func initPackage(name string, location *model.Location, manifestFile *zip.File, pomPropertiesFile *zip.File, dir *string, dockerTemp *string) error {
	endOfFile := regexp.MustCompile(jarPackagesRegex)
	pkg := new(model.Package)
	pkg.Metadata = Metadata{}
	pkg.ID = uuid.NewString()
	pkg.Path = name
	pkg.Type = java
	paths := strings.Split(location.Path, string(os.PathSeparator))
	pkg.Locations = append(pkg.Locations, model.Location{
		Path:      paths[len(paths)-1],
		LayerHash: location.LayerHash,
	})
	splitName := strings.Split(name, "/")
	fileName := splitName[len(splitName)-1]
	if manifestFile != nil {
		_ = parseManifest(manifestFile, pkg, dockerTemp)
		parseLicenses(pkg)
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
		parsePomProperties(string(data), pkg, pomPropertiesFile.Name)
		if err := reader.Close(); err != nil {
			return err
		}
	}
	if pomPropertiesFile != nil {
		vendor = pkg.Metadata.(Metadata)["PomProperties"]["artifactId"]
		version = pkg.Metadata.(Metadata)["PomProperties"]["version"]
		product = vendor
	} else {
		version = formatVersionMetadata(pkg, pkg.Metadata.(Metadata)["Manifest"]["Implementation-Version"])
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

	pkg.Name = vendor
	pkg.Version = version
	pkg.Description = pkg.Metadata.(Metadata)["Manifest"]["Bundle-Description"]

	parseJavaURL(pkg)

	cpe.NewCPE23(pkg, vendor, product, version)
	// additional CPEs if
	if pkg.Metadata.(Metadata)["PomProperties"] != nil {
		vendor = pkg.Metadata.(Metadata)["PomProperties"]["groupId"]
		if len(vendor) > 0 {
			generateAdditionalCPE(vendor, product, version, pkg)
		}
	}
	if pkg.Metadata.(Metadata)["Manifest"] != nil {
		vendor = pkg.Metadata.(Metadata)["Manifest"]["Automatic-Module-Name"]
		if len(vendor) > 0 {
			cpe.NewCPE23(pkg, vendor, product, version)
		}
		vendor = pkg.Metadata.(Metadata)["Manifest"]["Bundle-SymbolicName"]
		generateAdditionalCPE(vendor, product, version, pkg)
	}

	if len(*dir) > 0 || pkg.Name != "" && pkg.Version != "" {
		checkPackage(pkg, location.LayerHash)
	}
	return nil
}

func checkPackage(pkg *model.Package, layerHash string) {
	if _, exists := Result[pkg.Name+":"+pkg.Version+":"+layerHash]; !exists {
		Result[pkg.Name+":"+pkg.Version+":"+layerHash] = *pkg
	} else {
		_tmpPackage := Result[pkg.Name+":"+pkg.Version+":"+layerHash]
		if pkg.Metadata.(Metadata)["Manifest"] != nil {
			_tmpPackage.Metadata.(Metadata)["Manifest"] = pkg.Metadata.(Metadata)["Manifest"]
			_tmpPackage.CPEs = append(_tmpPackage.CPEs, pkg.CPEs...)
		}
		if pkg.Metadata.(Metadata)["PomProperties"] != nil {
			_tmpPackage.Metadata.(Metadata)["PomProperties"] = pkg.Metadata.(Metadata)["PomProperties"]
			_tmpPackage.CPEs = append(_tmpPackage.CPEs, pkg.CPEs...)
		}
		if pkg.Metadata.(Metadata)["PomProject"] != nil {
			_tmpPackage.Metadata.(Metadata)["PomProject"] = pkg.Metadata.(Metadata)["PomProject"]
			_tmpPackage.CPEs = append(_tmpPackage.CPEs, pkg.CPEs...)

		}
		_tmpPackage.CPEs = cpe.RemoveDuplicateCPES(_tmpPackage.CPEs)
		Result[pkg.Name+":"+pkg.Version+":"+layerHash] = _tmpPackage
	}

}

// Parse jar files
func parseJarFiles(dependencies []*zip.File, location *model.Location, dir *string, dockerTemp *string) error {

	if err := os.Mkdir(*dockerTemp, fs.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}

	for _, dependency := range dependencies {

		splitName := strings.Split(dependency.Name, "/")
		fileName := splitName[len(splitName)-1]
		jarFile, err := os.Create(filepath.Join(*dockerTemp, fileName))

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

		err = findManifestAndPomPropertiesFromDependencyJarFile(jarFile, location, dependency.Name, dir, dockerTemp)
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
func generateAdditionalCPE(vendor string, product string, version string, pkg *model.Package) {
	if len(vendor) > 0 {
		if strings.Contains(vendor, ".") {
			for _, v := range strings.Split(vendor, ".") {
				tldsRegex := `(com|org|io|edu|net|edu|gov|mil|the\ |a\ |an\ )(?:\b|')`
				if !regexp.MustCompile(tldsRegex).MatchString(v) {
					cpe.NewCPE23(pkg, v, product, version)
				}
			}
		} else {
			cpe.NewCPE23(pkg, vendor, product, version)
		}
	}
}

// Find the manifest and pom properties files from the dependency jar file
func findManifestAndPomPropertiesFromDependencyJarFile(jarFile *os.File, location *model.Location, name string, dir *string, dockerTemp *string) error {

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
		// Skip unsafe files fo extraction
		if strings.Contains(zipFile.Name, "..") {
			continue
		}
		if match := regexp.MustCompile(jarPackagesRegex).FindString(zipFile.Name); len(match) > 0 {
			dependencies = append(dependencies, zipFile)
		} else if strings.Contains(zipFile.Name, pomFileName) {
			file, err := os.Create(filepath.Join(*dockerTemp, pomFileName))
			if err != nil {
				panic(err)
			}
			reader, err := zipFile.Open()
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(file, reader)
			if err != nil {
				return err
			}

			paths := strings.Split(jarFile.Name(), "/")
			_ = parsePomXML(model.Location{
				Path:      file.Name(),
				LayerHash: location.LayerHash,
			}, paths[len(paths)-1], dir)

			if err := file.Close(); err != nil {
				return err
			}
			defer os.Remove(file.Name())
		}

		metadataFile, pomPropertiesFile = checkZipfile(zipFile, metadataFile, pomPropertiesFile)

		if metadataFile != nil || pomPropertiesFile != nil {
			_ = initPackage(name, location, metadataFile, pomPropertiesFile, dir, dockerTemp)
			pomPropertiesFile = nil
			metadataFile = nil
		}
	}
	_ = parseJarFiles(dependencies, location, dir, dockerTemp)
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
func parseLicenses(pkg *model.Package) {
	var licenses = make([]string, 0)
	for key, value := range pkg.Metadata.(Metadata)["Manifest"] {
		if strings.Contains(key, "Bundle-License") {
			licenses = append(licenses, strings.TrimSpace(value))
		}
	}

	pkg.Licenses = licenses
}

// Parse java manifest file
func parseManifest(manifestFile *zip.File, pkg *model.Package, dockerTemp *string) error {

	createdManifest, err := os.Create(filepath.Join(*dockerTemp, strings.Replace(manifestFile.Name, "/", "_", -1)))
	if err != nil {
		return err
	}
	defer createdManifest.Close()

	var value string
	var attribute string
	var previousAttribute string

	manifest := make(Manifest)
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
	pkg.Metadata.(Metadata)["ManifestLocation"] = Manifest{"path": filepath.Join(pkg.Path, manifestFile.Name)}
	pkg.Metadata.(Metadata)["Manifest"] = manifest
	return nil
}

// Parse pom properties
func parsePomProperties(data string, pkg *model.Package, path string) {

	var value string
	var attribute string

	pomProperties := make(Manifest)
	pomProperties["location"] = filepath.Join(pkg.Name, path)
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

	pkg.Metadata.(Metadata)["PomProperties"] = pomProperties
}

// Parse PURL
func parseJavaURL(pkg *model.Package) {
	if pkg.Metadata.(Metadata)["PomProperties"] != nil {
		pkg.PURL = model.PURL("pkg" + ":" + "maven" + "/" + pkg.Metadata.(Metadata)["PomProperties"]["groupId"] + "/" + pkg.Metadata.(Metadata)["PomProperties"]["artifactId"] + "@" + pkg.Version)
	} else {
		pkg.PURL = model.PURL("pkg" + ":" + "maven" + "/" + pkg.Name + "/" + pkg.Name + "@" + pkg.Version)
	}
}

func parsePomXML(location model.Location, layerPath string, dir *string) error {
	file, err := os.ReadFile(location.Path)
	if err != nil {
		return err
	}
	if err = xml.Unmarshal(file, &JavaPomXML); err != nil {
		return err
	}

	if len(JavaPomXML.Dependencies) > 0 {
		if len(*dir) > 0 {
			for _, dep := range JavaPomXML.Dependencies {
				if dep.ArtifactID != "" && !strings.Contains(dep.Version, "$") {
					pkg := new(model.Package)
					pkg.Metadata = Metadata{}
					pkg.ID = uuid.NewString()
					pkg.Name = dep.ArtifactID
					pkg.Path = util.TrimUntilLayer(model.Location{
						Path:      layerPath,
						LayerHash: location.LayerHash,
					})
					pkg.Version = dep.Version
					pkg.Type = java
					paths := strings.Split(location.Path, string(os.PathSeparator))
					pkg.Locations = append(pkg.Locations, model.Location{
						Path:      paths[len(paths)-1],
						LayerHash: location.LayerHash,
					})
					pkg.Metadata.(Metadata)["ManifestLocation"] = Manifest{"path": pkg.Path}
					pkg.Metadata.(Metadata)["PomProject"] = Manifest{
						"name":    pkg.Name,
						"version": pkg.Version,
						"groupID": dep.GroupID,
					}
					parseJavaURL(pkg)
					cpe.NewCPE23(pkg, pkg.Name, pkg.Name, pkg.Version)
					generateAdditionalCPE(dep.GroupID, pkg.Name, pkg.Version, pkg)
					checkPackage(pkg, location.LayerHash)
				}
			}
		} else {
			if _, exists := Result[JavaPomXML.ArtifactID+":"+JavaPomXML.Version+":"+location.LayerHash]; exists {
				_tmpPackage := Result[JavaPomXML.ArtifactID+":"+JavaPomXML.Version+":"+location.LayerHash]
				_tmpPackage.Metadata.(Metadata)["PomProject"] = Manifest{
					"name":    JavaPomXML.Name,
					"version": JavaPomXML.Version,
					"groupID": JavaPomXML.GroupID,
				}
				cpe.NewCPE23(&_tmpPackage, JavaPomXML.ArtifactID, JavaPomXML.ArtifactID, JavaPomXML.Version)
				generateAdditionalCPE(JavaPomXML.GroupID, JavaPomXML.ArtifactID, JavaPomXML.Version, &_tmpPackage)
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
		p.Metadata.(Metadata)["Manifest"]["Implementation-Version"] = versionMetadata[0]

		newMetadata := strings.SplitN(match+versionMetadata[1], ":", 2)
		if len(newMetadata) > 1 {
			p.Metadata.(Metadata)["Manifest"][newMetadata[0]] = strings.TrimSpace(newMetadata[1])
		}

		return versionMetadata[0]
	}

	return version
}
