package alpmdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
	"github.com/google/uuid"
)

const (
	alpmdb  = "alpmdb"
	libalpm = "libalpm"
)

var (
	// /var/lib/pacman/local/<package-version>/files
	installedPackagesPath = filepath.Join("var", "lib", "pacman", "local")
)

type Manifest map[string]interface{}

// FindAlpinePackagesFromContent check for alpine-os files in the file contents
func FindAlpmdbPackagesFromContent(req *bom.ParserRequirements) {
	if util.ParserEnabled(alpmdb, req.Arguments.EnabledParsers) {
		for _, content := range *req.Contents {
			if strings.Contains(content.Path, installedPackagesPath) {
				fmt.Printf("Alpmdb file found!!!: %v", content.Path)
				os.Exit(1)
				if err := parseInstalledPackages(content.Path, content.LayerHash, req.Arguments.DisableFileListing, req.Result.Packages); err != nil {
					err = errors.New("alpmdb-parser: " + err.Error())
					*req.Errors = append(*req.Errors, err)
				}
			}
		}
	}
	defer req.WG.Done()
}

// Init alpmdb package
func initAlpmdbPackage(pkg *model.Package) {
	pkg.Metadata = map[string]string{}
	pkg.ID = uuid.NewString()
	pkg.Type = alpmdb
	pkg.Path = installedPackagesPath
}

// Parse installed packages metadata
func parseInstalledPackages(filename string, layer string, noFileListing *bool, pkgs *[]model.Package) error {
	// Check if the file is valid
	if !isValidFile(filename) {
		return fmt.Errorf("%s is not a valid file", filename)
	}

	// Extract the package version from the filename
	version, err := getVersionFromPath(filename)
	if err != nil {
		return err
	}

	// Open the file
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Parse the package metadata
	var pkg model.Package
	var manifest Manifest
	if err := json.NewDecoder(f).Decode(&manifest); err != nil {
		return err
	}

	pkg.ID = uuid.New().String()
	pkg.Name = manifest["pkgname"].(string)
	pkg.Version = version
	pkg.Path = filename
	pkg.Type = libalpm

	// Check if package is not empty before append
	if pkg.Name != "" && pkg.Version != "" {
		*pkgs = append(*pkgs, pkg)
	}

	initAlpmdbPackage(&pkg)

	// Add the package to the list
	*pkgs = append(*pkgs, pkg)

	return nil
}

func getVersionFromPath(path string) (string, error) {
	// Extract the package version from the path using a regular expression.
	// For example, if the path is "/var/lib/pacman/local/foo-1.2.3/files",
	// this regular expression might match "1.2.3".
	re := regexp.MustCompile(`.*\/([^-\/]*)-([\d.]*(-r\d+)?)\/.*`)
	matches := re.FindStringSubmatch(path)

	if len(matches) < 3 {
		return "", errors.New("could not parse package version from path: " + path)
	}

	version := matches[2]
	return version, nil
}

func isValidFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}
