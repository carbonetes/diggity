package json

import (
	"fmt"

	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/pkg/types"
)

var log = logger.GetLogger()

func DisplayResults(result types.SoftwareManifest) {
	json, err := helper.ToJSON(result)
	if err != nil {
		log.Fatal("Error converting to JSON")
	}
	fmt.Println(string(json))
}
