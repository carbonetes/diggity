package hackage

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

const (
	hackage       = "hackage"
	stackYaml     = "stack.yaml"
	stackYamlLock = "stack.yaml.lock"
	cabalFreeze   = "cabal.project.freeze"
	shaTag        = "sha256"
	revTag        = "rev"
	anyTag        = "any."
	constraints   = "constraints:"
)

var (
	stackConfig     metadata.StackConfig
	stackLockConfig metadata.StackLockConfig
)

// FindHackagePackagesFromContent checks for stack.yaml, stack.yaml.lock, and cabal.project.freeze files in the file contents
func FindHackagePackagesFromContent(req *bom.ParserRequirements) {
	if util.ParserEnabled(hackage, req.Arguments.EnabledParsers) {
		for _, content := range *req.Contents {
			if filepath.Base(content.Path) == stackYaml {
				if err := readStackContent(&content, req.SBOM.Packages); err != nil {
					err = errors.New("hackage-parser: " + err.Error())
					*req.Errors = append(*req.Errors, err)
				}
			}
			if filepath.Base(content.Path) == stackYamlLock {
				if err := readStackLockContent(&content, req.SBOM.Packages); err != nil {
					err = errors.New("hackage-parser: " + err.Error())
					*req.Errors = append(*req.Errors, err)
				}
			}
			if filepath.Base(content.Path) == cabalFreeze {
				if err := readCabalFreezeContent(&content, req.SBOM.Packages); err != nil {
					err = errors.New("hackage-parser: " + err.Error())
					*req.Errors = append(*req.Errors, err)
				}
			}
		}
	}
	defer req.WG.Done()
}

// Read stack.yaml contents
func readStackContent(location *model.Location, pkgs *[]model.Package) error {
	stackBytes, _ := os.ReadFile(location.Path)
	err := yaml.Unmarshal(stackBytes, &stackConfig)

	if err != nil {
		// Skip invalid extra deps
		if strings.Contains(err.Error(), "cannot unmarshal !!map into string") {
			return nil
		}
		return err
	}

	for _, dep := range stackConfig.ExtraDeps {
		if name, _, _, _, _ := parseExtraDep(dep); name != "" {
			*pkgs = append(*pkgs, *initHackagePackage(location, dep, ""))
		}
	}

	return nil
}

// Read stack.yaml.lock contents
func readStackLockContent(location *model.Location, pkgs *[]model.Package) error {
	stackBytes, _ := os.ReadFile(location.Path)
	err := yaml.Unmarshal(stackBytes, &stackLockConfig)

	if err != nil {
		// Skip invalid extra deps
		if strings.Contains(err.Error(), "cannot unmarshal !!map into string") {
			return nil
		}
		return err
	}

	// Get snapshot URL
	snapshot := stackLockConfig.Snapshots[0].(map[string]interface{})["completed"]
	url := snapshot.(map[string]interface{})["url"].(string)

	for _, dep := range stackLockConfig.Packages {
		if name, _, _, _, _ := parseExtraDep(dep.Original.Hackage); name != "" {
			*pkgs = append(*pkgs, *initHackagePackage(location, dep.Original.Hackage, url))
		}
	}

	return nil
}

// Read cabal.project.freeze contents
func readCabalFreezeContent(location *model.Location, pkgs *[]model.Package) error {
	file, err := os.Open(location.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		var pkg string

		// Find packages by the any. tag
		if strings.Contains(line, anyTag) {
			// Remove constraints field
			if strings.Contains(line, constraints) {
				pkg = strings.Replace(line, constraints, "", -1)
			} else {
				pkg = line
			}
			if nv := formatCabalPackage(pkg); nv != "" {
				*pkgs = append(*pkgs, *initHackagePackage(location, nv, ""))
			}
		}
	}

	return nil
}

// Init Hackage Package
func initHackagePackage(location *model.Location, dep string, url string) *model.Package {
	name, version, pkgHash, size, rev := parseExtraDep(dep)

	pkg := new(model.Package)
	pkg.ID = uuid.NewString()
	pkg.Name = name
	pkg.Version = version
	pkg.Path = pkg.Name
	pkg.Type = hackage
	pkg.Locations = append(pkg.Locations, model.Location{
		Path:      util.TrimUntilLayer(*location),
		LayerHash: location.LayerHash,
	})
	pkg.Licenses = []string{}

	// get purl
	parseHackageURL(pkg)

	// get CPEs
	cpe.NewCPE23(pkg, "", pkg.Name, pkg.Version)

	// fill metadata
	pkg.Metadata = metadata.HackageMetadata{
		Name:        name,
		Version:     version,
		PkgHash:     pkgHash,
		Size:        size,
		Revision:    rev,
		SnapshotURL: url,
	}

	return pkg
}

// Parse PURL
func parseHackageURL(pkg *model.Package) {
	pkg.PURL = model.PURL("pkg" + ":" + hackage + "/" + pkg.Name + "@" + pkg.Version)
}

// Parse Name, Version, PkgHash, Size, and Revision from extra-deps
func parseExtraDep(dep string) (name string, version string, pkgHash string, size string, rev string) {
	pkg := strings.Split(dep, "@")
	nv := strings.Split(pkg[0], "-")
	name = strings.Join(nv[0:len(nv)-1], "-")
	version = nv[len(nv)-1]

	if len(pkg) > 1 {
		// Parse pkgHash if sha256 is detected
		if strings.Contains(pkg[1], shaTag) {
			hs := strings.Split(pkg[1], ",")
			pkgHash = hs[0]
			size = hs[1]
		}
		// Parse revision if rev is detected
		if strings.Contains(pkg[1], revTag) {
			rev = pkg[1]
		}
	}

	return name, version, pkgHash, size, rev
}

// Format Name and Version for parsing
func formatCabalPackage(anyPkg string) string {
	pkg := strings.Replace(strings.TrimSpace(anyPkg), anyTag, "", -1)
	nv := strings.Replace(pkg, " ==", "-", -1)
	return strings.Replace(nv, ",", "", -1)
}
