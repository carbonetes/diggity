package cpe

import (
	_ "embed"
	"encoding/json"
	"slices"
	"strings"

	"github.com/carbonetes/diggity/internal/log"
	prmt "github.com/gitchander/permutation"
)

// Collection of common keywords that cannot be used on each package
var excluded = []string{"cli", "v2", "net", "crypto", "sync"}

var (
	//go:embed data/product.json
	productEmbed []byte
	productMap   = map[string]string{}
)

func init() {
	if err := json.Unmarshal(productEmbed, &productMap); err != nil {
		log.Debugf("failed to unmarshal JSON: %v", err)
	}
}

func productLookup(name string) string {
	keywords := makeKeywords(name)
	removeExcluded(&keywords)
	for _, k := range keywords {
		if product, ok := productMap[k]; ok {
			return product
		}
	}
	return name
}

func makeKeywords(name string) []string {
	var keywords []string
	keywords = append(keywords, name)
	if strings.Contains(name, "-") {
		keyword := strings.ReplaceAll(name, "-", "_")
		if !slices.Contains(keywords, keyword) {
			keywords = append(keywords, keyword)
		}
		name = strings.ReplaceAll(name, "_", "-")
		parts := strings.Split(name, "-")
		keywords = shuffleKeywordParts(parts, keywords, "-")
		if len(parts) > 2 {
			parts = parts[:len(parts)-1]
			keyword = strings.Join(parts, "-")
			if !slices.Contains(keywords, keyword) {
				keywords = append(keywords, strings.Join(parts, "-"))
			}
			keywords = shuffleKeywordParts(parts, keywords, "-")
		}
	}
	if strings.Contains(name, "_") {
		keyword := strings.ReplaceAll(name, "_", "-")
		if !slices.Contains(keywords, keyword) {
			keywords = append(keywords, keyword)
		}
		name = strings.ReplaceAll(name, "-", "_")
		parts := strings.Split(name, "_")
		// keywords = shuffleKeywordParts(parts, keywords, "_")
		if len(parts) > 2 {
			parts = parts[:len(parts)-1]
			keyword = strings.Join(parts, "_")
			if !slices.Contains(keywords, keyword) {
				keywords = append(keywords, keyword)
			}
			keywords = shuffleKeywordParts(parts, keywords, "_")
		}
	}
	return keywords
}

// Create a new set of keywords by shuffling the parts a the string
func shuffleKeywordParts(parts, keywords []string, separator string) []string {
	if len(parts) > 0 {
		result := prmt.New(prmt.StringSlice(parts))
		for result.Next() {
			newKeyword := strings.Join(parts, separator)
			if !slices.Contains(keywords, newKeyword) {
				keywords = append(keywords, newKeyword)
			}
		}
	}
	return keywords
}

// Filter out keywords that cannot be used
func removeExcluded(keywords *[]string) {
	for index, k := range *keywords {
		for _, e := range excluded {
			if k == e {
				*keywords = append((*keywords)[:index], (*keywords)[index+1:]...)
			}
		}

	}
}
