package json

import (
	"fmt"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
)

func DisplayResults(result *cyclonedx.BOM) {
	json, err := helper.ToJSON(result)
	if err != nil {
		log.Debug("Error converting to JSON")
	}
	fmt.Println(string(json))
}
