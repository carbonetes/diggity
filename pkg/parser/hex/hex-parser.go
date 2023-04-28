package hex

import (
	"bufio"
	"errors"
	"regexp"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
	"github.com/google/uuid"

	"os"

	"path/filepath"
)

const (
	rebarLock = "rebar.lock"
	mixLock   = "mix.lock"
	hex       = "hex"
)

// Metadata hex metadata
type Metadata metadata.HexMetadata

var rebarLockRegex = regexp.MustCompile(`[\[{<">},: \]\n]+`)
var mixLockRegex = regexp.MustCompile(`[%{}\n" ,:]+`)

// FindHexPackagesFromContent - find hex packages from content
func FindHexPackagesFromContent(req *bom.ParserRequirements) {
	if util.ParserEnabled(hex, req.Arguments.EnabledParsers) {
		for _, content := range *req.Contents {
			if filepath.Base(content.Path) == rebarLock {
				if err := parseHexRebarPacakges(&content, req.SBOM.Packages); err != nil {
					err = errors.New("hex-parser: " + err.Error())
					*req.Errors = append(*req.Errors, err)
				}
			}
			if filepath.Base(content.Path) == mixLock {
				if err := parseHexMixPackages(&content, req.SBOM.Packages); err != nil {
					err = errors.New("hex-parser: " + err.Error())
					*req.Errors = append(*req.Errors, err)
				}
			}
		}
	}
	defer req.WG.Done()
}

// Parse hex package metadata - rebar
func parseHexRebarPacakges(location *model.Location, pkgs *[]model.Package) error {
	file, err := os.Open(location.Path)
	if err != nil {
		return err
	}
	// rebarMetadata := make(map[string]*model.Package)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		keyValue := scanner.Text()
		pkg := new(model.Package)
		tokens := rebarLockRegex.Split(keyValue, -1)

		if len(tokens) == 7 {
			name, version := tokens[1], tokens[4]
			pkg.ID = uuid.NewString()
			pkg.Name = name
			pkg.Version = version
			pkg.Type = hex
			pkg.Path = name
			pkg.Locations = append(pkg.Locations, model.Location{
				Path:      util.TrimUntilLayer(*location),
				LayerHash: location.LayerHash,
			})
			cpe.NewCPE23(pkg, pkg.Name, pkg.Name, pkg.Version)
			parseHexPURL(pkg)
			pkg.Metadata = Metadata{
				Name:    name,
				Version: version,
			}

		}
		if pkg.Name != "" {
			*pkgs = append(*pkgs, *pkg)
		}

	}
	return nil
}

// Parse hex package metadata - mix
func parseHexMixPackages(location *model.Location, pkgs *[]model.Package) error {
	file, err := os.Open(location.Path)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		keyValue := scanner.Text()
		pkg := new(model.Package)
		tokens := mixLockRegex.Split(keyValue, -1)

		if len(tokens) < 6 {
			continue
		}
		name, version, hash, hashExt := tokens[1], tokens[4], tokens[5], tokens[len(tokens)-2]

		pkg.ID = uuid.NewString()
		pkg.Name = name
		pkg.Version = version
		pkg.Type = hex
		pkg.Path = name
		pkg.Locations = append(pkg.Locations, model.Location{
			Path:      util.TrimUntilLayer(*location),
			LayerHash: location.LayerHash,
		})
		cpe.NewCPE23(pkg, pkg.Name, pkg.Name, pkg.Version)
		parseHexPURL(pkg)
		pkg.Metadata = Metadata{
			Name:       name,
			Version:    version,
			PkgHash:    hash,
			PkgHashExt: hashExt,
		}
		if pkg.Name != "" {
			*pkgs = append(*pkgs, *pkg)
		}
	}
	return nil
}

// Parse PURL
func parseHexPURL(pkg *model.Package) {
	pkg.PURL = model.PURL("pkg" + ":" + "hex" + "/" + pkg.Name + "@" + pkg.Version)
}
