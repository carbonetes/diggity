package portage

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/file"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
	"github.com/google/uuid"
)

var (
	portageDBPath    = filepath.Join("var", "db", "pkg") + string(os.PathSeparator)
	portage          = "portage"
	portageContent   = "CONTENTS"
	portageLicense   = "LICENSE"
	portageSize      = "SIZE"
	portageObj       = "obj"
	portageAlgorithm = "md5"
	ebuild           = "ebuild"
	noFileErrWin     = "The system cannot find the file specified"
	noFileErrMac     = "no such file or directory"
)

// FindPortagePackagesFromContent find portage metadata files
func FindPortagePackagesFromContent() {
	if util.ParserEnabled(portage) {
		for _, content := range file.Contents {
			if strings.Contains(content.Path, portageDBPath) && strings.Contains(content.Path, portageContent) {
				if err := readPortageContent(content); err != nil {
					err = errors.New("portage-parser: " + err.Error())
					bom.Errors = append(bom.Errors, &err)
				}
			}
		}
	}
	defer bom.WG.Done()
}

// Read Portage Contents
func readPortageContent(location *model.Location) error {
	// Parse package metadata from path
	pkg, err := initPortagePackage(location)
	if err != nil {
		return err
	}

	bom.Packages = append(bom.Packages, pkg)

	return nil
}

// Init Portage Package
func initPortagePackage(location *model.Location) (*model.Package, error) {
	_package := new(model.Package)
	_package.ID = uuid.NewString()

	contentPath := filepath.Dir(location.Path)
	name, version := portageNameVersion(contentPath)
	_package.Name = name
	_package.Version = version
	_package.Type = portage
	_package.Locations = append(_package.Locations, model.Location{
		Path:      util.TrimUntilLayer(*location),
		LayerHash: location.LayerHash,
	})
	_package.Path = strings.Split(contentPath, portageDBPath)[1]

	// get licenses
	if err := getPortageLicenses(_package, location.Path); err != nil {
		return _package, err
	}

	// get purl
	parsePortagePURL(_package)

	// get CPEs
	cpe.NewCPE23(_package, "", _package.Name, _package.Version)

	// fill metadata
	if err := initPortageMetadata(_package, location.Path); err != nil {
		return _package, err
	}

	return _package, nil
}

func initPortageMetadata(p *model.Package, loc string) error {
	var metadata metadata.PortageMetadata
	sizePath := strings.Replace(loc, portageContent, portageSize, -1)

	// Find and parse SIZE file
	file, err := os.Open(sizePath)
	if err != nil {
		if strings.Contains(err.Error(), noFileErrWin) || strings.Contains(err.Error(), noFileErrMac) {
			return nil
		}
		return err
	}

	// Get Size Metadata
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		size, err := strconv.Atoi(scanner.Text())
		if err != nil {
			continue
		}
		metadata.Size = size
	}

	// Get files metadata
	if !*bom.Arguments.DisableFileListing {
		if err := getPortageFiles(&metadata, loc); err != nil {
			return nil
		}
	}

	p.Metadata = metadata
	return nil
}

// Parse Portage Name and Version
func portageNameVersion(pkg string) (name string, version string) {
	// parse version
	r := regexp.MustCompile(`[0-9].*`)
	pkgBase := filepath.Base(pkg)
	version = r.FindString(pkgBase)

	// parse name
	namePath := strings.Split(pkg, portageDBPath)[1]
	name = strings.Replace(namePath, "-"+version, "", -1)

	return name, version
}

// Get Portage Licenses
func getPortageLicenses(p *model.Package, loc string) error {
	licenses := []string{}
	licensePath := strings.Replace(loc, portageContent, portageLicense, -1)

	// Find and parse LICENSE file
	file, err := os.Open(licensePath)
	if err != nil {
		if strings.Contains(err.Error(), noFileErrWin) || strings.Contains(err.Error(), noFileErrMac) {
			p.Licenses = licenses
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		licenses = append(licenses, scanner.Text())
	}
	p.Licenses = licenses

	return nil
}

// Get Portage Files
func getPortageFiles(md *metadata.PortageMetadata, loc string) error {
	var files []metadata.PortageFile

	// Parse CONTENT file
	file, err := os.Open(loc)
	if err != nil {
		if strings.Contains(err.Error(), noFileErrWin) || strings.Contains(err.Error(), noFileErrMac) {
			return nil
		}
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content := scanner.Text()
		if strings.Contains(content, portageObj) {
			files = append(files, parsePortageFile(content))
		}
	}

	md.Files = files

	return nil
}

// Parse Portage Files
func parsePortageFile(content string) metadata.PortageFile {
	var file metadata.PortageFile
	var digest metadata.PortageDigest

	obj := strings.Split(content, " ")
	// digest
	if len(obj) > 2 {
		digest.Algorithm = portageAlgorithm
		digest.Value = obj[2]
	}
	// file
	file.Path = obj[1]
	file.Digest = digest

	return file
}

// Parse PURL
func parsePortagePURL(_package *model.Package) {
	name := strings.Replace(_package.Name, string(os.PathSeparator), "/", -1)
	_package.PURL = model.PURL("pkg" + ":" + ebuild + "/" + name + "@" + _package.Version)
}
