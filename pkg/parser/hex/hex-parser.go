package parser

import (
	"bufio"
	"errors"
	"regexp"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/file"
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
func FindHexPackagesFromContent() {
	if util.ParserEnabled(hex) {
		for _, content := range file.Contents {
			if filepath.Base(content.Path) == rebarLock {
				if err := parseHexRebarPacakges(content); err != nil {
					err = errors.New("hex-parser: " + err.Error())
					bom.Errors = append(bom.Errors, &err)
				}
			}
			if filepath.Base(content.Path) == mixLock {
				if err := parseHexMixPackages(content); err != nil {
					err = errors.New("hex-parser: " + err.Error())
					bom.Errors = append(bom.Errors, &err)
				}
			}
		}
	}
	defer bom.WG.Done()
}

// Parse hex package metadata - rebar
func parseHexRebarPacakges(location *model.Location) error {
	file, err := os.Open(location.Path)
	if err != nil {
		return err
	}
	// rebarMetadata := make(map[string]*model.Package)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		keyValue := scanner.Text()
		_package := new(model.Package)
		tokens := rebarLockRegex.Split(keyValue, -1)

		if len(tokens) == 7 {
			name, version := tokens[1], tokens[4]
			_package.ID = uuid.NewString()
			_package.Name = name
			_package.Version = version
			_package.Type = hex
			_package.Path = name
			_package.Locations = append(_package.Locations, model.Location{
				Path:      util.TrimUntilLayer(*location),
				LayerHash: location.LayerHash,
			})
			cpe.NewCPE23(_package, _package.Name, _package.Name, _package.Version)
			parseHexPURL(_package)
			_package.Metadata = Metadata{
				Name:    name,
				Version: version,
			}

		}
		if _package.Name != "" {
			bom.Packages = append(bom.Packages, _package)
		}

	}
	return nil
}

// Parse hex package metadata - mix
func parseHexMixPackages(location *model.Location) error {
	file, err := os.Open(location.Path)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		keyValue := scanner.Text()
		_package := new(model.Package)
		tokens := mixLockRegex.Split(keyValue, -1)

		if len(tokens) < 6 {
			continue
		}
		name, version, hash, hashExt := tokens[1], tokens[4], tokens[5], tokens[len(tokens)-2]

		_package.ID = uuid.NewString()
		_package.Name = name
		_package.Version = version
		_package.Type = hex
		_package.Path = name
		_package.Locations = append(_package.Locations, model.Location{
			Path:      util.TrimUntilLayer(*location),
			LayerHash: location.LayerHash,
		})
		cpe.NewCPE23(_package, _package.Name, _package.Name, _package.Version)
		parseHexPURL(_package)
		_package.Metadata = Metadata{
			Name:       name,
			Version:    version,
			PkgHash:    hash,
			PkgHashExt: hashExt,
		}
		if _package.Name != "" {
			bom.Packages = append(bom.Packages, _package)
		}
	}
	return nil
}

// Parse PURL
func parseHexPURL(_package *model.Package) {
	_package.PURL = model.PURL("pkg" + ":" + "hex" + "/" + _package.Name + "@" + _package.Version)
}
