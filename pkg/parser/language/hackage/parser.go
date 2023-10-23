package hackage

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/common"
	"gopkg.in/yaml.v3"
)

const parserErr = "hackage-parser: "

var (
	stackConfig     metadata.StackConfig
	stackLockConfig metadata.StackLockConfig
)

// Read stack.yaml contents
func readStackContent(location *model.Location, req *common.ParserParams) {
	stackBytes, err := os.ReadFile(location.Path)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}
	err = yaml.Unmarshal(stackBytes, &stackConfig)

	if err != nil {
		// Skip invalid extra deps
		if strings.Contains(err.Error(), "cannot unmarshal !!map into string") {
			return
		}
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	for _, dep := range stackConfig.ExtraDeps {
		if name, _, _, _, _ := parseExtraDep(dep); name != "" {
			*req.SBOM.Packages = append(*req.SBOM.Packages, *initHackagePackage(location, dep, ""))
		}
	}
}

// Read stack.yaml.lock contents
func readStackLockContent(location *model.Location, req *common.ParserParams) {
	stackBytes, err := os.ReadFile(location.Path)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}
	err = yaml.Unmarshal(stackBytes, &stackLockConfig)

	if err != nil {
		// Skip invalid extra deps
		if strings.Contains(err.Error(), "cannot unmarshal !!map into string") {
			return
		}
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	// Get snapshot URL
	snapshot := stackLockConfig.Snapshots[0].(map[string]interface{})["completed"]
	url := snapshot.(map[string]interface{})["url"].(string)

	for _, dep := range stackLockConfig.Packages {
		if name, _, _, _, _ := parseExtraDep(dep.Original.Hackage); name != "" {
			*req.SBOM.Packages = append(*req.SBOM.Packages, *initHackagePackage(location, dep.Original.Hackage, url))
		}
	}
}

// Read cabal.project.freeze contents
func readCabalFreezeContent(location *model.Location, req *common.ParserParams) {
	file, err := os.Open(location.Path)
	if err != nil {
		err = errors.New(parserErr + err.Error())
		*req.Errors = append(*req.Errors, err)
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
				*req.SBOM.Packages = append(*req.SBOM.Packages, *initHackagePackage(location, nv, ""))
			}
		}
	}
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
