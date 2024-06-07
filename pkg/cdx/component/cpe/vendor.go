package cpe

import (
	_ "embed"
	"encoding/json"
	"log"
)

var (
	//go:embed data/vendor.json
	vendorEmbed []byte
	vendorMap   = map[string]string{}
)

func init() {
	// Unmarshal the JSON data into a map
	if err := json.Unmarshal(vendorEmbed, &vendorMap); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}
}

func vendorLookup(product string) (vendor string) {
	// Lookup the vendor in the map
	if vendor, ok := vendorMap[product]; ok {
		return vendor
	}

	return product
}
