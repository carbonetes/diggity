package cyclonedx

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/model/output"
	"github.com/carbonetes/diggity/internal/output/save"
	"github.com/carbonetes/diggity/internal/output/util"
	"github.com/carbonetes/diggity/internal/parser/bom"
	"github.com/carbonetes/diggity/internal/parser/distro"
	versionPackage "github.com/carbonetes/diggity/internal/version"

	"github.com/google/uuid"
)

const (
	vendor                                   = "carbonetes"
	name                                     = "diggity"
	diggityPrefix                            = "diggity"
	packagePrefix                            = "package"
	distroPrefix                             = "distro"
	colonPrefix                              = ":"
	cpePrefix                                = "cpe23"
	locationPrefix                           = "location"
	library          output.ComponentLibrary = "library"
	operatingSystem                          = "operating-system"
	issueTracker                             = "issue-tracker"
	referenceWebsite                         = "website"
	referenceOther                           = "other"
	// XMLN cyclonedx
	XMLN = "http://cyclonedx.org/schema/bom/1.4"
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

func convertPackage() *output.CycloneFormat {
	// Sort packages alphabetically
	util.SortPackages()

	//initialize component
	components := make([]output.Component, len(bom.Packages))
	for i, p := range bom.Packages {
		components[i] = convertToComponent(p)
	}

	components = append(components, addDistroComponent(distro.Distro()))

	return &output.CycloneFormat{
		XMLNS:        XMLN,
		SerialNumber: uuid.NewString(),
		Metadata:     getFromSource(),
		Components:   &components,
	}
}

func addDistroComponent(distro *model.Distro) output.Component {

	if distro == nil {
		return output.Component{}
	}
	externalReferences := &[]output.ExternalReference{}
	if distro.BugReportURL != "" {
		*externalReferences = append(*externalReferences, output.ExternalReference{
			URL:  distro.BugReportURL,
			Type: issueTracker,
		})
	}
	if distro.HomeURL != "" {
		*externalReferences = append(*externalReferences, output.ExternalReference{
			URL:  distro.HomeURL,
			Type: referenceWebsite,
		})
	}
	if distro.SupportURL != "" {
		*externalReferences = append(*externalReferences, output.ExternalReference{
			URL:     distro.SupportURL,
			Type:    referenceOther,
			Comment: "support",
		})
	}
	if distro.PrivacyPolicyURL != "" {
		*externalReferences = append(*externalReferences, output.ExternalReference{
			URL:     distro.PrivacyPolicyURL,
			Type:    referenceOther,
			Comment: "privacyPolicy",
		})
	}
	if len(*externalReferences) == 0 {
		externalReferences = nil
	}
	properties := make([]output.Property, 0)

	//assign id
	properties = append(properties, output.Property{
		Name:  diggityPrefix + colonPrefix + distroPrefix + ":id",
		Value: distro.ID,
	})

	properties = append(properties, output.Property{
		Name:  diggityPrefix + colonPrefix + distroPrefix + ":prettyName",
		Value: distro.PrettyName,
	})

	properties = append(properties, output.Property{
		Name:  diggityPrefix + colonPrefix + distroPrefix + ":distributionCodename",
		Value: distro.DistribCodename,
	})

	properties = append(properties, output.Property{
		Name:  diggityPrefix + colonPrefix + distroPrefix + ":versionID",
		Value: distro.VersionID,
	})

	return output.Component{
		Type:               operatingSystem,
		Name:               distro.ID,
		Description:        distro.PrettyName,
		ExternalReferences: externalReferences,
		Properties:         &properties,
	}
}

func getFromSource() *output.Metadata {
	//temp data-- data should come from final bom model
	versionInfo := versionPackage.FromBuild()
	return &output.Metadata{
		Timestamp: time.Now().Format(time.RFC3339),
		Tools: &[]output.Tool{
			{
				Vendor:  vendor,
				Name:    name,
				Version: versionInfo.Version,
			},
		},
	}
}

func convertToComponent(p *model.Package) output.Component {

	return output.Component{
		BOMRef:     addID(p),
		Type:       library,
		Name:       p.Name,
		Version:    p.Version,
		PackageURL: string(p.PURL),
		Licenses:   convertLicense(p),
		Properties: initProperties(p),
	}
}

func initProperties(p *model.Package) *[]output.Property {
	properties := make([]output.Property, 0)

	//assign type
	properties = append(properties, output.Property{
		Name:  diggityPrefix + colonPrefix + packagePrefix + ":type",
		Value: p.Type,
	})

	//assign cpes
	for _, cpe := range p.CPEs {
		properties = append(properties, output.Property{
			Name:  diggityPrefix + colonPrefix + cpePrefix,
			Value: cpe,
		})
	}

	//assign locations
	for i, location := range p.Locations {
		index := strconv.Itoa(i)

		//add hash
		properties = append(properties, output.Property{
			Name:  diggityPrefix + colonPrefix + locationPrefix + colonPrefix + index + colonPrefix + "layerHash",
			Value: location.LayerHash,
		})
		//add path
		properties = append(properties, output.Property{
			Name:  diggityPrefix + colonPrefix + locationPrefix + colonPrefix + index + colonPrefix + "path",
			Value: location.Path,
		})

	}

	return &properties
}

func addID(p *model.Package) string {
	return string(p.PURL) + "?package-id=" + p.ID
}

func convertLicense(p *model.Package) *[]output.License {
	// lm := output.LicenseModel{}
	licenses := make([]output.License, 0)
	for _, licenseName := range p.Licenses {
		licenses = append(licenses, output.License{
			ID: licenseName,
		})
	}
	if len(licenses) > 0 {
		return &licenses
	}
	return nil

}
