package alpmdb

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
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
			if strings.Contains(content.Path, installedPackagesPath) && strings.Contains(content.Path, "\\desc") {

				if err := readDesc(content.Path, req.SBOM.Packages, content.LayerHash); err != nil {
					err = errors.New("alpmdb-parser: " + err.Error())
					*req.Errors = append(*req.Errors, err)
				}
			}
		}
	}
	defer req.WG.Done()
}

func readDesc(path string, pkgs *[]model.Package, layer string) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	contents := string(data)
	metadata := make(map[string]string)
	lines := strings.Split(contents, "\n")
	for index, line := range lines {
		line = strings.TrimSpace(line)

		if line == "%PROVIDES%" {
			for i := index + 1; i < len(lines); i++ {
				l := strings.TrimSpace(lines[i])
				if !(l == "") || !(strings.HasPrefix(l, "%")) {
					metadata["provides"] += l + " "
				}
			}
		} else if !strings.HasPrefix(line, "%") {
			if len(line) > 0 {
				key, value := parseMetadataLine(line)
				if key != "" {
					metadata[key] = value
				}
			}
		} else {
			if strings.Contains(line, "%") {
				key := strings.ToLower(strings.ReplaceAll(line, "%", ""))
				value := strings.TrimSpace(lines[index+1])
				metadata[key] = value
			}
		}
	}

	pkg := newPackage(metadata, layer)

	generateCPE(&pkg)

	*pkgs = append(*pkgs, pkg)

	return nil
}

func parseMetadataLine(line string) (string, string) {
	fields := strings.SplitN(line, ":", 2)
	if len(fields) == 2 {
		return strings.TrimSpace(fields[0]), strings.TrimSpace(fields[1])
	}
	return "", ""
}

func newPackage(metadata map[string]string, layer string) model.Package {

	return model.Package{
		ID:          uuid.NewString(),
		Name:        metadata["name"],
		Type:        alpmdb,
		Version:     metadata["version"],
		Path:        installedPackagesPath,
		Locations:   generateLocations(layer),
		Description: metadata["desc"],
		Licenses:    generateLicenses(metadata["license"]),
		PURL:        generatePURL(metadata),
		Metadata:    metadata,
	}
}

func generateLocations(layer string) []model.Location {
	return []model.Location{
		{
			LayerHash: layer,
			Path:      installedPackagesPath,
		},
	}
}

func generateLicenses(value string) []string {
	var licenses []string
	for _, license := range strings.Split(value, " ") {
		if !strings.Contains(strings.ToLower(license), "and") {
			licenses = append(licenses, license)
		}
	}
	return licenses
}

func generatePURL(metadata map[string]string) model.PURL {
	arch, ok := metadata["arch"]
	if !ok {
		arch = ""
	}
	origin, ok := metadata["origin"]
	if !ok {
		origin = ""
	}
	return model.PURL("pkg" + `:` + alpmdb + `/` + "archlinux" + `/` + metadata["name"] + `@` + metadata["version"] + `?arch=` + arch + `&` + `upstream=` + origin + `&distro=` + "archlinux")
}

func generateCPE(pkg *model.Package) {
	cpe.NewCPE23(pkg, "", pkg.Name, pkg.Version)
}
