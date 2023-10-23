package maven

import (
	"archive/zip"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/language/java/gradle"
	"github.com/google/uuid"
)

//TODO: reduce the code complexities here

// Init java package
func initPackage(name string, location *model.Location, manifestFile *zip.File, pomPropertiesFile *zip.File, dir *string, dockerTemp *string, result *map[string]model.Package) error {
	endOfFile := regexp.MustCompile(jarPackagesRegex)
	pkg := new(model.Package)
	pkg.Metadata = Metadata{}
	pkg.ID = uuid.NewString()
	pkg.Path = name
	pkg.Type = Type
	pkg.PackageOrigin = model.ApplicationPackage
	pkg.Parser = Type
	pkg.Language = gradle.Language
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
		if strings.EqualFold(version, "") {
			version = "0.0.0"
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
		checkPackage(pkg, location.LayerHash, result)
	}

	return nil
}

func checkPackage(pkg *model.Package, layerHash string, result *map[string]model.Package) {
	if _, exists := (*result)[pkg.Name+":"+pkg.Version+":"+layerHash]; !exists {
		(*result)[pkg.Name+":"+pkg.Version+":"+layerHash] = *pkg
	} else {
		tmpPackage := (*result)[pkg.Name+":"+pkg.Version+":"+layerHash]
		if pkg.Metadata.(Metadata)["Manifest"] != nil {
			tmpPackage.Metadata.(Metadata)["Manifest"] = pkg.Metadata.(Metadata)["Manifest"]
			tmpPackage.CPEs = append(tmpPackage.CPEs, pkg.CPEs...)
		}
		if pkg.Metadata.(Metadata)["PomProperties"] != nil {
			tmpPackage.Metadata.(Metadata)["PomProperties"] = pkg.Metadata.(Metadata)["PomProperties"]
			tmpPackage.CPEs = append(tmpPackage.CPEs, pkg.CPEs...)
		}
		if pkg.Metadata.(Metadata)["PomProject"] != nil {
			tmpPackage.Metadata.(Metadata)["PomProject"] = pkg.Metadata.(Metadata)["PomProject"]
			tmpPackage.CPEs = append(tmpPackage.CPEs, pkg.CPEs...)

		}
		tmpPackage.CPEs = cpe.RemoveDuplicateCPES(tmpPackage.CPEs)
		(*result)[pkg.Name+":"+pkg.Version+":"+layerHash] = tmpPackage
	}

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

// Parse PURL
func parseJavaURL(pkg *model.Package) {
	if pkg.Metadata.(Metadata)["PomProperties"] != nil {
		pkg.PURL = model.PURL("pkg" + ":" + "maven" + "/" + pkg.Metadata.(Metadata)["PomProperties"]["groupId"] + "/" + pkg.Metadata.(Metadata)["PomProperties"]["artifactId"] + "@" + pkg.Version)
	} else {
		pkg.PURL = model.PURL("pkg" + ":" + "maven" + "/" + pkg.Name + "/" + pkg.Name + "@" + pkg.Version)
	}
}
