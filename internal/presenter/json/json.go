package json

import (
	"fmt"

	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/types"
)

func DisplayResults(result types.SoftwareManifest) {
	json, err := helper.ToJSON(result)
	if err != nil {
		log.Error("Error converting to JSON")
	}
	fmt.Println(string(json))
}
