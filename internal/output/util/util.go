package util

import (
	"sort"

	"github.com/carbonetes/diggity/pkg/parser/bom"
)

// SortPackages sort packages by name alphabetically
func SortPackages() {
	sort.Slice(bom.Packages, func(i, j int) bool {
		return bom.Packages[i].Name < bom.Packages[j].Name
	})
}
