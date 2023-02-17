package cpe

import (
	"strings"

	"github.com/carbonetes/diggity/internal/model"
)

// NewCPE23 Generates and Validates CPE String based on CPE Version 2.3
func NewCPE23(_package *model.Package, vendor string, product string, version string) *model.Package {
	baseCPE := toCPE(vendor, product, version)
	if _package.Type == "java" && strings.Contains(baseCPE.Vendor, ";") {
		for _, _vendor := range strings.Split(baseCPE.Vendor, ";") {
			baseCPE.Vendor = _vendor
			_package.CPEs = append(_package.CPEs, expandCPEsBySeparators(*baseCPE)...)
		}
	} else {
		_package.CPEs = append(_package.CPEs, cpeToString(*baseCPE))
		_package.CPEs = append(_package.CPEs, expandCPEsBySeparators(*baseCPE)...)
	}

	baseCPE.Vendor = baseCPE.Product
	_package.CPEs = append(_package.CPEs, cpeToString(*baseCPE))
	_package.CPEs = RemoveDuplicateCPES(_package.CPEs)

	// Retain base CPE
	if len(_package.CPEs) == 0 {
		_package.CPEs = append(_package.CPEs, cpeToString(*baseCPE))
	}

	return _package
}

func cpeJoin(matchers ...string) string {
	var cpe string
	for i, matcher := range matchers {
		if i > 0 {
			cpe += ":" + matcher
		} else {
			cpe += matcher
		}
	}
	return cpe
}

func cpeToString(baseCPE CPE) string {
	return cpeJoin("cpe:2.3", baseCPE.Part, baseCPE.Vendor, baseCPE.Product, baseCPE.Version, baseCPE.Update, baseCPE.Edition, baseCPE.Language, baseCPE.SWEdition, baseCPE.TargetSW, baseCPE.TargetHW, baseCPE.Other)
}

func toCPE(vendor string, product string, version string) *CPE {

	return &CPE{
		Part:      "a",
		Vendor:    vendor,
		Product:   product,
		Version:   version,
		Update:    wildcard,
		Edition:   wildcard,
		SWEdition: wildcard,
		TargetSW:  wildcard,
		TargetHW:  wildcard,
		Other:     wildcard,
		Language:  wildcard,
	}
}

// Separators = "-", "_", " ", "."
func expandCPEsBySeparators(baseCPE CPE) []string {

	cpes := make([]string, 0)

	// Vendor
	if strings.Contains(baseCPE.Vendor, "-") || strings.Contains(baseCPE.Product, "-") {
		tmp := baseCPE
		for _, vendor := range expand(baseCPE, "Vendor", '-', '_') {
			baseCPE.Vendor = vendor
			cpes = append(cpes, cpeToString(baseCPE))
			for _, product := range expand(baseCPE, "Product", '-', '_') {
				baseCPE.Product = product
				cpes = append(cpes, cpeToString(baseCPE))
			}
		}

		// Return to original value
		baseCPE = tmp
	} else if strings.Contains(baseCPE.Vendor, "_") || strings.Contains(baseCPE.Product, "_") {
		tmp := baseCPE
		for _, vendor := range expand(baseCPE, "Vendor", '_', '-') {
			baseCPE.Vendor = vendor
			cpes = append(cpes, cpeToString(baseCPE))
			for _, product := range expand(baseCPE, "Product", '_', '-') {
				baseCPE.Product = product
				cpes = append(cpes, cpeToString(baseCPE))
			}
		}

		// Return to original value
		baseCPE = tmp
	} else if strings.Contains(baseCPE.Vendor, ".") {
		tmp := baseCPE
		for _, vendor := range strings.Split(baseCPE.Vendor, ".") {
			baseCPE.Vendor = vendor
			cpes = append(cpes, cpeToString(baseCPE))
		}

		// Return to original value
		baseCPE = tmp
	}

	return cpes
}

func expand(baseCPE CPE, _field field, separator rune, replace rune) []string {
	expandedFields := make([]string, 0)
	switch _field {
	case "Vendor":
		{
			vendorBytes := []byte(baseCPE.Vendor)
			expandedFields = append(expandedFields, string(vendorBytes))
			for idx, c := range baseCPE.Vendor {
				if c == separator {
					vendorBytes[idx] = byte(replace)
					expandedFields = append(expandedFields, string(vendorBytes))
				}
				if c == replace {
					vendorBytes[idx] = byte(separator)
					expandedFields = append(expandedFields, string(vendorBytes))
				}
			}

			return expandedFields
		}
	case "Product":
		{
			productBytes := []byte(baseCPE.Product)
			expandedFields = append(expandedFields, string(productBytes))
			for idx, c := range baseCPE.Product {
				if c == separator {
					productBytes[idx] = byte(replace)
					expandedFields = append(expandedFields, string(productBytes))
				}
				if c == replace {
					productBytes[idx] = byte(separator)
					expandedFields = append(expandedFields, string(productBytes))
				}
			}
			return expandedFields
		}
	}

	return nil
}
