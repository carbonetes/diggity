package maven

import (
	"archive/zip"
	"bufio"
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

//TODO: reduce the code complexities here

// Extract jar files
func extractJarFile(location *model.Location, dir *string, dockerTemp *string, result *map[string]model.Package) error {

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
			}, location.Path, dir, result)

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
			_ = initPackage(paths[len(paths)-1], location, metadataFile, pomPropertiesFile, dir, dockerTemp, result)
		}
	}

	err = parseJarFiles(dependencies, location, dir, dockerTemp, result)
	return err
}

// Parse jar files
func parseJarFiles(dependencies []*zip.File, location *model.Location, dir *string, dockerTemp *string, result *map[string]model.Package) error {

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

		err = findManifestAndPomPropertiesFromDependencyJarFile(jarFile, location, dependency.Name, dir, dockerTemp, result)
		if err != nil {
			return err
		}

		if err := jarFile.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Find the manifest and pom properties files from the dependency jar file
func findManifestAndPomPropertiesFromDependencyJarFile(jarFile *os.File, location *model.Location, name string, dir *string, dockerTemp *string, result *map[string]model.Package) error {

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
			}, paths[len(paths)-1], dir, result)

			if err := file.Close(); err != nil {
				return err
			}
			defer os.Remove(file.Name())
		}

		metadataFile, pomPropertiesFile = checkZipfile(zipFile, metadataFile, pomPropertiesFile)

		if metadataFile != nil || pomPropertiesFile != nil {
			_ = initPackage(name, location, metadataFile, pomPropertiesFile, dir, dockerTemp, result)
			pomPropertiesFile = nil
			metadataFile = nil
		}
	}
	_ = parseJarFiles(dependencies, location, dir, dockerTemp, result)
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
