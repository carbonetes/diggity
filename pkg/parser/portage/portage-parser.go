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
func FindPortagePackagesFromContent(req *bom.ParserRequirements) {
	if util.ParserEnabled(portage, req.Arguments.EnabledParsers) {
		for _, content := range *req.Contents {
			if strings.Contains(content.Path, portageDBPath) && strings.Contains(content.Path, portageContent) {
				if err := readPortageContent(&content, req.Arguments.DisableFileListing, req.Result.Packages); err != nil {
					err = errors.New("portage-parser: " + err.Error())
					*req.Errors = append(*req.Errors, err)
				}
			}
		}
	}
	defer req.WG.Done()
}

// Read Portage Contents
func readPortageContent(location *model.Location, noFileListing *bool, pkgs *[]model.Package) error {
	// Parse package metadata from path
	pkg, err := initPortagePackage(location, noFileListing)
	if err != nil {
		return err
	}

	*pkgs = append(*pkgs, *pkg)

	return nil
}

// Init Portage Package
func initPortagePackage(location *model.Location, noFileListing *bool) (*model.Package, error) {
	pkg := new(model.Package)
	pkg.ID = uuid.NewString()

	contentPath := filepath.Dir(location.Path)
	name, version := portageNameVersion(contentPath)
	pkg.Name = name
	pkg.Version = version
	pkg.Type = portage
	pkg.Locations = append(pkg.Locations, model.Location{
		Path:      util.TrimUntilLayer(*location),
		LayerHash: location.LayerHash,
	})
	pkg.Path = strings.Split(contentPath, portageDBPath)[1]

	// get licenses
	if err := getPortageLicenses(pkg, location.Path); err != nil {
		return pkg, err
	}

	// get purl
	parsePortagePURL(pkg)

	// get CPEs
	cpe.NewCPE23(pkg, "", pkg.Name, pkg.Version)

	// fill metadata
	if err := initPortageMetadata(pkg, location.Path, noFileListing); err != nil {
		return pkg, err
	}

	return pkg, nil
}

func initPortageMetadata(p *model.Package, loc string, noFileListing *bool) error {
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
	if !*noFileListing {
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
func parsePortagePURL(pkg *model.Package) {
	name := strings.Replace(pkg.Name, string(os.PathSeparator), "/", -1)
	pkg.PURL = model.PURL("pkg" + ":" + ebuild + "/" + name + "@" + pkg.Version)
}
