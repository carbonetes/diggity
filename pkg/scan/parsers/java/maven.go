package java

import (
	"encoding/xml"
	"io/ioutil"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan/parsers/common"
)

// MavenParser handles parsing of pom.xml files
type MavenParser struct {
	base *common.BaseParser
}

// NewMavenParser creates a new Maven parser
func NewMavenParser(base *common.BaseParser) *MavenParser {
	return &MavenParser{base: base}
}

// MavenProject represents the structure of a pom.xml file
type MavenProject struct {
	XMLName      xml.Name          `xml:"project"`
	GroupID      string            `xml:"groupId"`
	ArtifactID   string            `xml:"artifactId"`
	Version      string            `xml:"version"`
	Dependencies MavenDependencies `xml:"dependencies"`
}

// MavenDependencies represents the dependencies section
type MavenDependencies struct {
	Dependency []MavenDependency `xml:"dependency"`
}

// MavenDependency represents a single dependency
type MavenDependency struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
	Scope      string `xml:"scope"`
}

// Parse parses pom.xml files
func (m *MavenParser) Parse(filePath string) ([]model.Package, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var project MavenProject
	if err := xml.Unmarshal(content, &project); err != nil {
		return nil, err
	}

	var packages []model.Package

	for _, dep := range project.Dependencies.Dependency {
		if dep.GroupID == "" || dep.ArtifactID == "" {
			continue
		}

		// Create package name as groupId:artifactId
		name := dep.GroupID + ":" + dep.ArtifactID

		// Check if it's a dev dependency
		isDev := dep.Scope == "test" || dep.Scope == "provided"

		metadata := map[string]interface{}{
			"manager":     "maven",
			"group_id":    dep.GroupID,
			"artifact_id": dep.ArtifactID,
			"scope":       dep.Scope,
		}

		pkg := m.base.CreatePackage(name, dep.Version, filePath, isDev, metadata)
		packages = append(packages, pkg)
	}

	return packages, nil
}
