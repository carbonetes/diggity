package ci

import (
	"github.com/CycloneDX/cyclonedx-go"
)

type Table struct {
	Component *cyclonedx.Component
	Type      string
}

type ComponentTally struct {
	ComponentType string
	Count         int
}

func Evaluate(cdx *cyclonedx.BOM) []ComponentTally {
	tallyMap := make(map[string]int)
	for _, c := range *cdx.Components {
		for _, p := range *c.Properties {
			if p.Name == "diggity:package:type" {
				tallyMap[string(p.Value)]++
			}
		}
	}
	var tally []ComponentTally
	for name, count := range tallyMap {
		tally = append(tally, ComponentTally{
			ComponentType: name,
			Count:         count,
		})
	}
	return tally
}
