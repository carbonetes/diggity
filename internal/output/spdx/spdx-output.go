package spdx

import (
	"fmt"

	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/internal/output/json.go"
	"github.com/carbonetes/diggity/internal/output/save"
	"github.com/carbonetes/diggity/pkg/convert"
	"github.com/carbonetes/diggity/pkg/model"
	"gopkg.in/yaml.v3"
)

var log = logger.GetLogger()

// PrintSpdxJSON Print Packages in SPDX-JSON format
func PrintSpdxJSON(args *model.Arguments, outputType *string, pkgs *[]model.Package) {
	spdx := convert.ToSPDX(args, pkgs)
	result, err := json.ToJSON(spdx)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		panic(err)
	}

	if len(*args.OutputFile) > 0 {
		save.ResultToFile(string(result), outputType, args.OutputFile)
	} else {
		fmt.Printf("%+v\n", string(result))
	}
}

// PrintSpdxTagValue Print Packages in SPDX-TAG_VALUE format
func PrintSpdxTagValue(args *model.Arguments, outputType *string, pkgs *[]model.Package) {
	spdx := convert.ToSPDXTagValue(args, pkgs)

	if len(*args.OutputFile) > 0 {
		save.ResultToFile(stringSliceToString(*spdx),outputType, args.OutputFile)
	} else {
		fmt.Printf("%+v", stringSliceToString(*spdx))
	}
}

// PrintSpdxYaml Print Packages in SPDX Yaml format
func PrintSpdxYaml(args *model.Arguments, outputType *string, pkgs *[]model.Package) {
	spdx := convert.ToSPDX(args, pkgs)
	result, err := yaml.Marshal(spdx)
	if err != nil {
		log.Fatal(err)
	}

	if len(*args.OutputFile) > 0 {
		save.ResultToFile(string(result), outputType, args.OutputFile)
	} else {
		fmt.Printf("%+v\n", string(result))
	}
}

// convert spdx-tag-values to single string
func stringSliceToString(slice []string) string {
	result := ""
	for _, s := range slice {
		result += fmt.Sprintln(s)
	}
	return result
}
