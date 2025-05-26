package ci

import (
	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/types"
)

func Run(cdx *cyclonedx.BOM, secrets interface{}) {

	// BOM
	totalComponents := len(*cdx.Components)

	log.Printf("\nTotal Packages Found: %v", totalComponents)
	if totalComponents > 0 {
		result := Evaluate(cdx)
		TallyTable(result)
		log.Print("\nPackages Found")
		ResultTable(cdx)
	}

	// Secrets
	secretList := secrets.([]types.Secret)
	totalSecrets := len(secretList)
	log.Printf("\nTotal Secrets Found: %v", totalSecrets)
	if totalSecrets > 0 {
		SecretTable(secretList)
		log.Print("\nAssessment Failed: Exposed Secrets!")
	}

}
