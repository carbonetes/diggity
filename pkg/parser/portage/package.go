package portage

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/util"
	"github.com/google/uuid"
)

// Init Portage Package
func initPortagePackage(location *model.Location, noFileListing *bool) (*model.Package, error) {
	pkg := new(model.Package)
	pkg.ID = uuid.NewString()

	contentPath := filepath.Dir(location.Path)
	name, version := portageNameVersion(contentPath)
	pkg.Name = name
	pkg.Version = version
	pkg.Type = Type
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

// Parse PURL
func parsePortagePURL(pkg *model.Package) {
	name := strings.Replace(pkg.Name, string(os.PathSeparator), "/", -1)
	pkg.PURL = model.PURL("pkg" + ":" + ebuild + "/" + name + "@" + pkg.Version)
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
