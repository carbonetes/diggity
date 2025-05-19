package ci

import (
	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/log"
)

func Run(cdx *cyclonedx.BOM) {
	totalComponents := len(*cdx.Components)

	log.Printf("\nTotal Packages Found: %v", totalComponents)
	result := Evaluate(cdx)
	TallyTable(result)
	log.Print("\nPackages Found")
	ResultTable(cdx)

}
