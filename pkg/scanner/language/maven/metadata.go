package maven

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"strings"

	"golang.org/x/net/html/charset"
)

type POM struct {
	Path   string
	Files  []string
	Parent string
}

type Metadata struct {
	XMLName                xml.Name                `xml:"project" json:"project,omitempty"`
	ModelVersion           string                  `xml:"modelVersion" json:"modelVersion,omitempty"`
	Parent                 Parent                  `xml:"parent" json:"parent,omitempty"`
	GroupID                string                  `xml:"groupId" json:"groupId,omitempty"`
	ArtifactID             string                  `xml:"artifactId" json:"artifactId,omitempty"`
	Version                string                  `xml:"version" json:"version,omitempty"`
	Packaging              string                  `xml:"packaging" json:"packaging,omitempty"`
	Name                   string                  `xml:"name" json:"name,omitempty"`
	Description            string                  `xml:"description" json:"description,omitempty"`
	URL                    string                  `xml:"url" json:"url,omitempty"`
	InceptionYear          string                  `xml:"inceptionYear" json:"inceptionYear,omitempty"`
	Organization           *Organization           `xml:"organization" json:"organization,omitempty"`
	Licenses               *[]License              `xml:"licenses>license" json:"licenses,omitempty"`
	Developers             *[]Developer            `xml:"developers>developer" json:"developers,omitempty"`
	Contributors           *[]Contributor          `xml:"contributors>contributor" json:"contributors,omitempty"`
	MailingLists           *[]MailingList          `xml:"mailingLists>mailingList" json:"mailingLists,omitempty"`
	Prerequisites          *Prerequisites          `xml:"prerequisites" json:"prerequisites,omitempty"`
	Modules                []string                `xml:"modules>module" json:"modules,omitempty"`
	SCM                    *Scm                    `xml:"scm" json:"scm,omitempty"`
	IssueManagement        *IssueManagement        `xml:"issueManagement" json:"issueManagement,omitempty"`
	CIManagement           *CIManagement           `xml:"ciManagement" json:"ciManagement,omitempty"`
	DistributionManagement *DistributionManagement `xml:"distributionManagement" json:"distributionManagement,omitempty"`
	DependencyManagement   *DependencyManagement   `xml:"dependencyManagement" json:"dependencyManagement,omitempty"`
	Dependencies           []Dependency            `xml:"dependencies>dependency" json:"dependencies,omitempty"`
	Repositories           []PomRepository         `xml:"repositories>repository" json:"repositories,omitempty"`
	PluginRepositories     []PluginRepository      `xml:"pluginRepositories>pluginRepository" json:"pluginRepositories,omitempty"`
	Build                  *Build                  `xml:"build" json:"build,omitempty"`
	Reporting              *Reporting              `xml:"reporting" json:"reporting,omitempty"`
	Profiles               *[]Profile              `xml:"profiles>profile" json:"profiles,omitempty"`
	Properties             *Properties             `xml:"properties" json:"properties,omitempty"`
}

type Property struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

type Properties struct {
	Properties []Property `xml:",any"`
}

type Parent struct {
	GroupID      string `xml:"groupId" json:"groupId,omitempty"`
	ArtifactID   string `xml:"artifactId" json:"artifactId,omitempty"`
	Version      string `xml:"version" json:"version,omitempty"`
	RelativePath string `xml:"relativePath" json:"relativePath,omitempty"`
}

type Organization struct {
	Name string `xml:"name" json:"name,omitempty"`
	URL  string `xml:"url" json:"url,omitempty"`
}

type License struct {
	Name         string `xml:"name" json:"name,omitempty"`
	URL          string `xml:"url" json:"url,omitempty"`
	Distribution string `xml:"distribution" json:"distribution,omitempty"`
	Comments     string `xml:"comments" json:"comments,omitempty"`
}

type Developer struct {
	ID              string      `xml:"id" json:"id,omitempty"`
	Name            string      `xml:"name" json:"name,omitempty"`
	Email           string      `xml:"email" json:"email,omitempty"`
	URL             string      `xml:"url" json:"url,omitempty"`
	Organization    string      `xml:"organization" json:"organization,omitempty"`
	OrganizationURL string      `xml:"organizationUrl" json:"organizationUrl,omitempty"`
	Roles           []string    `xml:"roles>role" json:"roles,omitempty"`
	Timezone        string      `xml:"timezone" json:"timezone,omitempty"`
	Properties      *Properties `xml:"properties" json:"properties,omitempty"`
}

type Contributor struct {
	Name            string      `xml:"name" json:"name,omitempty"`
	Email           string      `xml:"email" json:"email,omitempty"`
	URL             string      `xml:"url" json:"url,omitempty"`
	Organization    string      `xml:"organization" json:"organization,omitempty"`
	OrganizationURL string      `xml:"organizationUrl" json:"organizationUrl,omitempty"`
	Roles           []string    `xml:"roles>role" json:"roles,omitempty"`
	Timezone        string      `xml:"timezone" json:"timezone,omitempty"`
	Properties      *Properties `xml:"properties" json:"properties,omitempty"`
}

type MailingList struct {
	Name          string   `xml:"name" json:"name,omitempty"`
	Subscribe     string   `xml:"subscribe" json:"subscribe,omitempty"`
	Unsubscribe   string   `xml:"unsubscribe" json:"unsubscribe,omitempty"`
	Post          string   `xml:"post" json:"post,omitempty"`
	Archive       string   `xml:"archive" json:"archive,omitempty"`
	OtherArchives []string `xml:"otherArchives>otherArchive" json:"otherArchives,omitempty"`
}

type Prerequisites struct {
	Maven string `xml:"maven" json:"maven,omitempty"`
}

type Scm struct {
	Connection          string `xml:"connection" json:"connection,omitempty"`
	DeveloperConnection string `xml:"developerConnection" json:"developerConnection,omitempty"`
	Tag                 string `xml:"tag" json:"tag,omitempty"`
	URL                 string `xml:"url" json:"url,omitempty"`
}

type IssueManagement struct {
	System string `xml:"system" json:"system,omitempty"`
	URL    string `xml:"url" json:"url,omitempty"`
}

type CIManagement struct {
	System    string     `xml:"system" json:"system,omitempty"`
	URL       string     `xml:"url" json:"url,omitempty"`
	Notifiers []Notifier `xml:"notifiers>notifier" json:"notifiers,omitempty"`
}

type Notifier struct {
	Type          string      `xml:"type" json:"type,omitempty"`
	SendOnError   bool        `xml:"sendOnError" json:"sendOnError,omitempty"`
	SendOnFailure bool        `xml:"sendOnFailure" json:"sendOnFailure,omitempty"`
	SendOnSuccess bool        `xml:"sendOnSuccess" json:"sendOnSuccess,omitempty"`
	SendOnWarning bool        `xml:"sendOnWarning" json:"sendOnWarning,omitempty"`
	Address       string      `xml:"address" json:"address,omitempty"`
	Configuration *Properties `xml:"configuration" json:"configuration,omitempty"`
}

type Repository struct {
	Type string `json:"type,omitempty" mapstructure:"type"`
	URL  string `json:"url,omitempty" mapstructure:"url"`
}

type DistributionManagement struct {
	Repository         *Repository `xml:"repository" json:"repository,omitempty"`
	SnapshotRepository *Repository `xml:"snapshotRepository" json:"snapshotRepository,omitempty"`
	Site               *Site       `xml:"site" json:"site,omitempty"`
	DownloadURL        string      `xml:"downloadUrl" json:"downloadUrl,omitempty"`
	Relocation         *Relocation `xml:"relocation" json:"relocation,omitempty"`
	Status             string      `xml:"status" json:"status,omitempty"`
}

type Site struct {
	ID   string `xml:"id" json:"id,omitempty"`
	Name string `xml:"name" json:"name,omitempty"`
	URL  string `xml:"url" json:"url,omitempty"`
}

type Relocation struct {
	GroupID    string `xml:"groupId" json:"groupId,omitempty"`
	ArtifactID string `xml:"artifactId" json:"artifactId,omitempty"`
	Version    string `xml:"version" json:"version,omitempty"`
	Message    string `xml:"message" json:"message,omitempty"`
}

type DependencyManagement struct {
	Dependencies []Dependency `xml:"dependencies>dependency" json:"dependencies,omitempty"`
}

type Dependency struct {
	GroupID    string      `xml:"groupId" json:"groupId,omitempty"`
	ArtifactID string      `xml:"artifactId" json:"artifactId,omitempty"`
	Version    string      `xml:"version" json:"version,omitempty"`
	Type       string      `xml:"type" json:"type,omitempty"`
	Classifier string      `xml:"classifier" json:"classifier,omitempty"`
	Scope      string      `xml:"scope" json:"scope,omitempty"`
	SystemPath string      `xml:"systemPath" json:"systemPath,omitempty"`
	Exclusions []Exclusion `xml:"exclusions>exclusion" json:"exclusions,omitempty"`
	Optional   string      `xml:"optional" json:"optional,omitempty"`
}

type Exclusion struct {
	ArtifactID string `xml:"artifactId" json:"artifactId,omitempty"`
	GroupID    string `xml:"groupId" json:"groupId,omitempty"`
}

type PomRepository struct {
	UniqueVersion bool              `xml:"uniqueVersion" json:"uniqueVersion,omitempty"`
	Releases      *RepositoryPolicy `xml:"releases" json:"releases,omitempty"`
	Snapshots     *RepositoryPolicy `xml:"snapshots" json:"snapshots,omitempty"`
	ID            string            `xml:"id" json:"id,omitempty"`
	Name          string            `xml:"name" json:"name,omitempty"`
	URL           string            `xml:"url" json:"url,omitempty"`
	Layout        string            `xml:"layout" json:"layout,omitempty"`
}

type RepositoryPolicy struct {
	Enabled        string `xml:"enabled" json:"enabled,omitempty"`
	UpdatePolicy   string `xml:"updatePolicy" json:"updatePolicy,omitempty"`
	ChecksumPolicy string `xml:"checksumPolicy" json:"checksumPolicy,omitempty"`
}

type PluginRepository struct {
	Releases  *RepositoryPolicy `xml:"releases" json:"releases,omitempty"`
	Snapshots *RepositoryPolicy `xml:"snapshots" json:"snapshots,omitempty"`
	ID        string            `xml:"id" json:"id,omitempty"`
	Name      string            `xml:"name" json:"name,omitempty"`
	URL       string            `xml:"url" json:"url,omitempty"`
	Layout    string            `xml:"layout" json:"layout,omitempty"`
}

type BuildBase struct {
	DefaultGoal      string           `xml:"defaultGoal" json:"defaultGoal,omitempty"`
	Resources        []Resource       `xml:"resources>resource" json:"resources,omitempty"`
	TestResources    []Resource       `xml:"testResources>testResource" json:"testResources,omitempty"`
	Directory        string           `xml:"directory" json:"directory,omitempty"`
	FinalName        string           `xml:"finalName" json:"finalName,omitempty"`
	Filters          []string         `xml:"filters>filter" json:"filters,omitempty"`
	PluginManagement PluginManagement `xml:"pluginManagement" json:"pluginManagement,omitempty"`
	Plugins          []Plugin         `xml:"plugins>plugin" json:"plugins,omitempty"`
}

type Build struct {
	SourceDirectory       string      `xml:"sourceDirectory" json:"sourceDirectory,omitempty"`
	ScriptSourceDirectory string      `xml:"scriptSourceDirectory" json:"scriptSourceDirectory,omitempty"`
	TestSourceDirectory   string      `xml:"testSourceDirectory" json:"testSourceDirectory,omitempty"`
	OutputDirectory       string      `xml:"outputDirectory" json:"outputDirectory,omitempty"`
	TestOutputDirectory   string      `xml:"testOutputDirectory" json:"testOutputDirectory,omitempty"`
	Extensions            []Extension `xml:"extensions>extension" json:"extensions,omitempty"`
	BuildBase
}

type Extension struct {
	GroupID    string `xml:"groupId" json:"groupId,omitempty"`
	ArtifactID string `xml:"artifactId" json:"artifactId,omitempty"`
	Version    string `xml:"version" json:"version,omitempty"`
}

type Resource struct {
	TargetPath string   `xml:"targetPath" json:"targetPath,omitempty"`
	Filtering  string   `xml:"filtering" json:"filtering,omitempty"`
	Directory  string   `xml:"directory" json:"directory,omitempty"`
	Includes   []string `xml:"includes>include" json:"includes,omitempty"`
	Excludes   []string `xml:"excludes>exclude" json:"excludes,omitempty"`
}

type PluginManagement struct {
	Plugins []Plugin `xml:"plugins>plugin" json:"plugins,omitempty"`
}

type Plugin struct {
	GroupID      string            `xml:"groupId" json:"groupId,omitempty"`
	ArtifactID   string            `xml:"artifactId" json:"artifactId,omitempty"`
	Version      string            `xml:"version" json:"version,omitempty"`
	Extensions   string            `xml:"extensions" json:"extensions,omitempty"`
	Executions   []PluginExecution `xml:"executions>execution" json:"executions,omitempty"`
	Dependencies []Dependency      `xml:"dependencies>dependency" json:"dependencies,omitempty"`
	Inherited    string            `xml:"inherited" json:"inherited,omitempty"`
}

type PluginExecution struct {
	ID        string   `xml:"id" json:"id,omitempty"`
	Phase     string   `xml:"phase" json:"phase,omitempty"`
	Goals     []string `xml:"goals>goal" json:"goals,omitempty"`
	Inherited string   `xml:"inherited" json:"inherited,omitempty"`
}

type Reporting struct {
	ExcludeDefaults string            `xml:"excludeDefaults" json:"excludeDefaults,omitempty"`
	OutputDirectory string            `xml:"outputDirectory" json:"outputDirectory,omitempty"`
	Plugins         []ReportingPlugin `xml:"plugins>plugin" json:"plugins,omitempty"`
}

type ReportingPlugin struct {
	GroupID    string      `xml:"groupId" json:"groupId,omitempty"`
	ArtifactID string      `xml:"artifactId" json:"artifactId,omitempty"`
	Version    string      `xml:"version" json:"version,omitempty"`
	Inherited  string      `xml:"inherited" json:"inherited,omitempty"`
	ReportSets []ReportSet `xml:"reportSets>reportSet" json:"reportSets,omitempty"`
}

type ReportSet struct {
	ID        string   `xml:"id" json:"id,omitempty"`
	Reports   []string `xml:"reports>report" json:"reports,omitempty"`
	Inherited string   `xml:"inherited" json:"inherited,omitempty"`
}

type Profile struct {
	ID                     string                  `xml:"id" json:"id,omitempty"`
	Activation             *Activation             `xml:"activation" json:"activation,omitempty"`
	Build                  *BuildBase              `xml:"build" json:"build,omitempty"`
	Modules                []string                `xml:"modules>module" json:"modules,omitempty"`
	DistributionManagement *DistributionManagement `xml:"distributionManagement" json:"distributionManagement,omitempty"`
	Properties             *Properties             `xml:"properties" json:"properties,omitempty"`
	DependencyManagement   *DependencyManagement   `xml:"dependencyManagement" json:"dependencyManagement,omitempty"`
	Dependencies           []Dependency            `xml:"dependencies>dependency" json:"dependencies,omitempty"`
	Repositories           []Repository            `xml:"repositories>repository" json:"repositories,omitempty"`
	PluginRepositories     []PluginRepository      `xml:"pluginRepositories>pluginRepository" json:"pluginRepositories,omitempty"`
	Reporting              *Reporting              `xml:"reporting" json:"reporting,omitempty"`
}

type Activation struct {
	ActiveByDefault bool                `xml:"activeByDefault" json:"activeByDefault,omitempty"`
	JDK             string              `xml:"jdk" json:"jdk,omitempty"`
	OS              *ActivationOS       `xml:"os" json:"os,omitempty"`
	Property        *ActivationProperty `xml:"property" json:"property,omitempty"`
	File            *ActivationFile     `xml:"file" json:"file,omitempty"`
}

type ActivationOS struct {
	Name    string `xml:"name" json:"name,omitempty"`
	Family  string `xml:"family" json:"family,omitempty"`
	Arch    string `xml:"arch" json:"arch,omitempty"`
	Version string `xml:"version" json:"version,omitempty"`
}

type ActivationProperty struct {
	Name  string `xml:"name" json:"name,omitempty"`
	Value string `xml:"value" json:"value,omitempty"`
}

type ActivationFile struct {
	Missing string `xml:"missing" json:"missing,omitempty"`
	Exists  string `xml:"exists" json:"exists,omitempty"`
}

// parsePOM parses the POM file and returns the metadata.
//
//nolint:all
func parsePOM(content []byte) (*Metadata, error) {
	// Create a new XML decoder which can handle various charsets
	decoder := xml.NewDecoder(bytes.NewReader(content))
	decoder.CharsetReader = charset.NewReaderLabel

	// Decode the XML into the Metadata struct
	var metadata Metadata
	err := decoder.Decode(&metadata)
	if err != nil {
		if err == io.EOF {
			return nil, errors.New("empty POM file")
		}
		return nil, err
	}

	return &metadata, nil
}

// resolveProperties resolves the properties in the POM file.
//
//nolint:all
func resolveProperties(metadata *Metadata) map[string]string {
	properties := make(map[string]string)
	if metadata.Properties != nil {
		for _, prop := range metadata.Properties.Properties {
			properties["${"+prop.XMLName.Local+"}"] = prop.Value
		}
	}
	return properties
}

// parseManifestFile parses the MANIFEST.MF file and returns a map of key-value pairs.
//
//nolint:all
func parseManifestFile(content []byte) (map[string]string, error) {
	manifest := make(map[string]string)

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			manifest[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return manifest, nil
}

// parsePOMProperties parses the pom.properties file and returns a map of key-value pairs.
func parsePOMProperties(content []byte) (map[string]string, error) {
	properties := make(map[string]string)

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			properties[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return properties, nil
}
