package cyclonedx

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/internal/output/save"
	"github.com/carbonetes/diggity/pkg/convert"
	"github.com/carbonetes/diggity/pkg/model"
)

var log = logger.GetLogger()

// PrintCycloneDXXML Print Packages in XML format
func PrintCycloneDXXML(pkgs *[]model.Package, outputType *string, filename *string) {
	cdx := convert.ToCDX(pkgs)
	result, err := xml.MarshalIndent(cdx, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	if len(*filename) > 0 {
		save.ResultToFile(string(result), outputType, filename)
	} else {
		fmt.Printf("%+v\n", string(result))
	}
}

// PrintCycloneDXJSON Print Packages in Cyclonedx Json format
func PrintCycloneDXJSON(pkgs *[]model.Package, outputType *string, filename *string) {
	cdx := convert.ToCDX(pkgs)
	result, err := json.MarshalIndent(cdx, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	if len(*filename) > 0 {
		save.ResultToFile(string(result), outputType, filename)
	} else {
		fmt.Printf("%+v\n", string(result))
	}
}
