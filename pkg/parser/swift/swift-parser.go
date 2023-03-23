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
		_package := new(model.Package)
		_package.ID = uuid.NewString()

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

		_package.Name = name
		_package.Version = version
		_package.Type = pod
		_package.Path = name
		_package.Locations = append(_package.Locations, model.Location{
			Path:      util.TrimUntilLayer(*location),
			LayerHash: location.LayerHash,
		})
		basepodname := strings.Split(name, "/")[0]
		cpe.NewCPE23(_package, name, name, version)
		parseSwiftPURL(_package)
		if val, ok := podFileLockFileMetadata.SpecChecksums[basepodname]; ok {
			_package.Metadata = metadata.PodFileLockMetadataCheckSums{
				Checksums: val,
			}
		} else {
			return nil
		}

		bom.Packages = append(bom.Packages, _package)
	}

	return nil
}

// Parse PURL
func parseSwiftPURL(_package *model.Package) {
	_package.PURL = model.PURL("pkg" + ":" + "cocoapods" + "/" + _package.Name + "@" + _package.Version)
}
