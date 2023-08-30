package pnpm

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
	"gopkg.in/yaml.v3"
)

const parserErr string = "pnpm-parser: "

var packageNameRegex = regexp.MustCompile(`^/?([^(]*)(?:\(.*\))*$`)

func readLockFile(location *model.Location, req *bom.ParserRequirements) {
	var lockfile metadata.PnpmMetadata

	file, err := os.ReadFile(location.Path)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
		return
	}

	err = yaml.Unmarshal(file, &lockfile)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
		return
	}

	for name, value := range lockfile.Dependencies {
		var version string

		if len(name) == 0 {
			continue
		}

		switch value := value.(type) {
		case map[string]interface{}:
			v, ok := value["version"].(string)
			if !ok {
				break
			}
			version = strings.SplitN(v, "(", 2)[0]
		case string:
			version = strings.SplitN(value, "(", 2)[0]
		default:
			break
		}

		if len(version) == 0 {
			continue
		}

		pkg := newPackage(name, version)
		generateCPEs(pkg)
		pkg.Path = util.TrimUntilLayer(*location)
		pkg.Locations = append(pkg.Locations, model.Location{
			Path:      pkg.Path,
			LayerHash: location.LayerHash,
		})
		pkg.Metadata = value
		*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
	}

	separator := "/"
	lockfileVersion, err := strconv.ParseFloat(lockfile.LockFileVersion, 64)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
		return
	}

	if lockfileVersion >= 6.0 {
		separator = "@"
	}

	for value, meta := range lockfile.Packages {
		value = packageNameRegex.ReplaceAllString(value, "$1")
		value = strings.TrimPrefix(value, "/")
		values := strings.Split(value, separator)

		name := strings.Join(values[:len(values)-1], separator)
		version := values[len(values)-1]

		pkg := newPackage(name, version)
		generateCPEs(pkg)
		pkg.Path = util.TrimUntilLayer(*location)
		pkg.Locations = append(pkg.Locations, model.Location{
			Path:      pkg.Path,
			LayerHash: location.LayerHash,
		})
		pkg.Metadata = meta
		*req.SBOM.Packages = append(*req.SBOM.Packages, *pkg)
	}
}
