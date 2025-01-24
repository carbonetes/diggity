package ui

import (
	"fmt"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
)

func printJson(result *cyclonedx.BOM) {
	json, err := helper.ToJSON(result)
	if err != nil {
		log.Debug("Error converting to JSON")
	}
	fmt.Println(string(json))
}
