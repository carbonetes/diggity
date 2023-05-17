package cyclonedx

import (
	"encoding/xml"
	"fmt"

	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/internal/output/json.go"
	"github.com/carbonetes/diggity/internal/output/save"
	"github.com/carbonetes/diggity/pkg/convert"
	"github.com/carbonetes/diggity/pkg/model"
)

var log = logger.GetLogger()

// PrintCycloneDXXML Print Packages in XML format
func PrintCycloneDXXML(sbom *model.SBOM, filename *string) {
	cdx := convert.ToCDX(sbom)
	result, err := xml.MarshalIndent(cdx, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	if len(*filename) > 0 {
		save.ResultToFile(string(result), filename)
	} else {
		fmt.Printf("%+v\n", string(result))
	}
}

// PrintCycloneDXJSON Print Packages in Cyclonedx Json format
func PrintCycloneDXJSON(sbom *model.SBOM, filename *string) {
	cdx := convert.ToCDX(sbom)
	result, err := json.ToJSON(cdx)
	if err != nil {
		log.Fatal(err)
	}
	if len(*filename) > 0 {
		save.ResultToFile(string(result), filename)
	} else {
		fmt.Printf("%+v\n", string(result))
	}
}
