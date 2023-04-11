package cyclonedx

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	"github.com/carbonetes/diggity/internal/output/save"
	"github.com/carbonetes/diggity/internal/output/util"
	versionPackage "github.com/carbonetes/diggity/internal/version"
	"github.com/carbonetes/diggity/pkg/model"

	cyclonedx "github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/distro"

	"github.com/google/uuid"
)

const (
	cycloneDX        = "CycloneDX"
	vendor           = "carbonetes"
	name             = "diggity"
	diggityPrefix    = "diggity"
	packagePrefix    = "package"
	distroPrefix     = "distro"
	colonPrefix      = ":"
	cpePrefix        = "cpe23"
	locationPrefix   = "location"
	library          = "library"
	operatingSystem  = "operating-system"
	issueTracker     = "issue-tracker"
	referenceWebsite = "website"
	referenceOther   = "other"
	version          = 1
)

var (
	// XMLN cyclonedx
	XMLN = fmt.Sprintf("http://cyclonedx.org/schema/bom/%+v", cyclonedx.SpecVersion1_4)
)

// PrintCycloneDXXML Print Packages in XML format
func PrintCycloneDXXML() {

	cyclonedxOuput := convertPackage()

	result, _ := xml.MarshalIndent(cyclonedxOuput, "", " ")
	if len(*bom.Arguments.OutputFile) > 0 {
		save.ResultToFile(string(result))
	} else {
		fmt.Printf("%+v\n", string(result))
	}
}

// PrintCycloneDXJSON Print Packages in Cyclonedx Json format
func PrintCycloneDXJSON() {

	cyclonedxOuput := convertPackage()

	result, _ := json.MarshalIndent(cyclonedxOuput, "", " ")

	if len(*bom.Arguments.OutputFile) > 0 {
		save.ResultToFile(string(result))
	} else {
		fmt.Printf("%+v\n", string(result))
	}
}

func convertPackage() *cyclonedx.BOM {
	// Sort packages alphabetically
	util.SortPackages()

	//initialize component
	components := make([]cyclonedx.Component, len(bom.Packages))
	for i, p := range bom.Packages {
		components[i] = convertToComponent(p)
	}

	components = append(components, addDistroComponent(distro.Distro()))

	return &cyclonedx.BOM{
		BOMFormat:    cycloneDX,
		SpecVersion:  cyclonedx.SpecVersion1_4,
		XMLNS:        XMLN,
		SerialNumber: uuid.NewString(),
		Version:      version,
		Metadata:     getFromSource(),
		Components:   &components,
	}
}

func addDistroComponent(distro *model.Distro) cyclonedx.Component {

	if distro == nil {
		return cyclonedx.Component{}
	}
	externalReferences := &[]cyclonedx.ExternalReference{}
	if distro.BugReportURL != "" {
		*externalReferences = append(*externalReferences, cyclonedx.ExternalReference{
			URL:  distro.BugReportURL,
			Type: issueTracker,
		})
	}
	if distro.HomeURL != "" {
		*externalReferences = append(*externalReferences, cyclonedx.ExternalReference{
			URL:  distro.HomeURL,
			Type: referenceWebsite,
		})
	}
	if distro.SupportURL != "" {
		*externalReferences = append(*externalReferences, cyclonedx.ExternalReference{
			URL:     distro.SupportURL,
			Type:    referenceOther,
			Comment: "support",
		})
	}
	if distro.PrivacyPolicyURL != "" {
		*externalReferences = append(*externalReferences, cyclonedx.ExternalReference{
			URL:     distro.PrivacyPolicyURL,
			Type:    referenceOther,
			Comment: "privacyPolicy",
		})
	}
	if len(*externalReferences) == 0 {
		externalReferences = nil
	}
	properties := make([]cyclonedx.Property, 0)

	//assign id
	properties = append(properties, cyclonedx.Property{
		Name:  diggityPrefix + colonPrefix + distroPrefix + ":id",
		Value: distro.ID,
	})

	properties = append(properties, cyclonedx.Property{
		Name:  diggityPrefix + colonPrefix + distroPrefix + ":prettyName",
		Value: distro.PrettyName,
	})

	properties = append(properties, cyclonedx.Property{
		Name:  diggityPrefix + colonPrefix + distroPrefix + ":distributionCodename",
		Value: distro.DistribCodename,
	})

	properties = append(properties, cyclonedx.Property{
		Name:  diggityPrefix + colonPrefix + distroPrefix + ":versionID",
		Value: distro.VersionID,
	})

	return cyclonedx.Component{
		Type:               operatingSystem,
		Name:               distro.ID,
		Description:        distro.PrettyName,
		ExternalReferences: externalReferences,
		Properties:         &properties,
	}
}

func getFromSource() *cyclonedx.Metadata {
	//temp data-- data should come from final bom model
	versionInfo := versionPackage.FromBuild()
	return &cyclonedx.Metadata{
		Timestamp: time.Now().Format(time.RFC3339),
		Tools: &[]cyclonedx.Tool{
			{
				Vendor:  vendor,
				Name:    name,
				Version: versionInfo.Version,
			},
		},
	}
}

func convertToComponent(p *model.Package) cyclonedx.Component {

	return cyclonedx.Component{
		BOMRef:     addID(p),
		Type:       library,
		Name:       p.Name,
		Version:    p.Version,
		PackageURL: string(p.PURL),
		Licenses:   convertLicense(p),
		Properties: initProperties(p),
	}
}

func initProperties(p *model.Package) *[]cyclonedx.Property {
	properties := make([]cyclonedx.Property, 0)

	//assign type
	properties = append(properties, cyclonedx.Property{
		Name:  diggityPrefix + colonPrefix + packagePrefix + ":type",
		Value: p.Type,
	})

	//assign cpes
	for _, cpe := range p.CPEs {
		properties = append(properties, cyclonedx.Property{
			Name:  diggityPrefix + colonPrefix + cpePrefix,
			Value: cpe,
		})
	}

	//assign locations
	for i, location := range p.Locations {
		index := strconv.Itoa(i)

		//add hash
		properties = append(properties, cyclonedx.Property{
			Name:  diggityPrefix + colonPrefix + locationPrefix + colonPrefix + index + colonPrefix + "layerHash",
			Value: location.LayerHash,
		})
		//add path
		properties = append(properties, cyclonedx.Property{
			Name:  diggityPrefix + colonPrefix + locationPrefix + colonPrefix + index + colonPrefix + "path",
			Value: location.Path,
		})

	}

	return &properties
}

func addID(p *model.Package) string {
	return string(p.PURL) + "?package-id=" + p.ID
}

func convertLicense(p *model.Package) *cyclonedx.Licenses {
	licenses := make(cyclonedx.Licenses, 0)

	// Get Licenses for CycloneDX model
	for _, licenseName := range p.Licenses {
		license := cyclonedx.License{
			ID: licenseName,
		}
		licenses = append(licenses, cyclonedx.LicenseChoice{
			License: &license,
		})
	}

	if len(licenses) > 0 {
		return &licenses
	}

	return &licenses
}
