package swift

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/internal/file"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

const (
	podfilelock = "Podfile.lock"
	pubname     = "name"
	pod         = "pod"
)

var podFileLockFileMetadata metadata.PodFileLockMetadata

// FindSwiftPackagesFromContent - find swift and objective-c packages from content
func FindSwiftPackagesFromContent() {
	if util.ParserEnabled(pod) {
		for _, content := range file.Contents {
			if filepath.Base(content.Path) == podfilelock {
				if err := parseSwiftPackages(content); err != nil {
					err = errors.New("swift-parser: " + err.Error())
					bom.Errors = append(bom.Errors, &err)
				}
			}
		}
	}
	defer bom.WG.Done()
}

func parseSwiftPackages(location *model.Location) error {
	byteValue, err := os.ReadFile(location.Path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal([]byte(byteValue), &podFileLockFileMetadata); err != nil {
		return err
	}

	for _, cPackage := range podFileLockFileMetadata.Pods {
		pkg := new(model.Package)
		pkg.ID = uuid.NewString()

		//check metadata in podfile lock file
		var pods string
		switch cp := cPackage.(type) {
		case string:
			pods = cp
		case map[string]interface{}:
			podVal := cPackage.(map[string]interface{})
			for podsAll := range podVal {
				pods = podsAll
			}
		}

		splits := strings.Split(pods, " ")
		name := splits[0]
		version := strings.TrimSuffix(strings.TrimPrefix(splits[1], "("), ")")

		pkg.Name = name
		pkg.Version = version
		pkg.Type = pod
		pkg.Path = name
		pkg.Locations = append(pkg.Locations, model.Location{
			Path:      util.TrimUntilLayer(*location),
			LayerHash: location.LayerHash,
		})
		basepodname := strings.Split(name, "/")[0]
		cpe.NewCPE23(pkg, name, name, version)
		parseSwiftPURL(pkg)
		if val, ok := podFileLockFileMetadata.SpecChecksums[basepodname]; ok {
			pkg.Metadata = metadata.PodFileLockMetadataCheckSums{
				Checksums: val,
			}
		} else {
			return nil
		}

		bom.Packages = append(bom.Packages, pkg)
	}

	return nil
}

// Parse PURL
func parseSwiftPURL(pkg *model.Package) {
	pkg.PURL = model.PURL("pkg" + ":" + "cocoapods" + "/" + pkg.Name + "@" + pkg.Version)
}
