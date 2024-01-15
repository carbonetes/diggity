package cpe

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/facebookincubator/nvdtools/wfn"
)

type (
	// CPE = wfn.Attributes
	CPE   = wfn.Attributes
	field = string
)

const (
	// Source: https://csrc.nist.gov/schema/cpe/2.3/cpe-naming_2.3.xsd
	cpeRegexString = `cpe:2\.3:[aho\*\-](:(((\?*|\*?)([a-zA-Z0-9\-\._]|(\\[\\\*\?!"#$$%&'\(\)\+,\/:;<=>@\[\]\^\x60\{\|}~]))+(\?*|\*?))|[\*\-]|[\+])){5}(:(([a-zA-Z]{2,3}(-([a-zA-Z]{2}|[0-9]{3}))?)|[\*\-]))(:(((\?*|\*?)([a-zA-Z0-9\-\._]|(\\[\\\*\?!"#$$%&'\(\)\+,\/:;<=>@\[\]\^\x60\{\|}~]))+(\?*|\*?))|[\*\-])){4}`
	wildcard       = "*"
)

var regExp = regexp.MustCompile(cpeRegexString)

// RemoveDuplicateCPES removes duplicate CPEs
func RemoveDuplicateCPES(cpes []string) []string {
	processed := make(map[string]bool)
	var list []string
	for _, cpe := range cpes {
		if _, value := processed[cpe]; !value {
			processed[cpe] = true
			if err := validateCPE(cpe); err == nil {
				list = append(list, cpe)
			}
		}
	}
	return list
}

func validateCPE(cpe string) error {
	if !regExp.MatchString(cpe) {
		return fmt.Errorf("failed to create CPE, invalid CPE string")
	}
	return nil
}

// NewCPE23 Generates and Validates CPE String based on CPE Version 2.3
func NewCPE23(vendor, product, version, category string) []string {
	var cpes []string
	baseCPE := toCPE(vendor, product, version)
	if category == "java" && strings.Contains(baseCPE.Vendor, ";") {
		for _, _vendor := range strings.Split(baseCPE.Vendor, ";") {
			baseCPE.Vendor = _vendor
			cpes = append(cpes, expandCPEsBySeparators(*baseCPE)...)
		}
	} else {
		cpes = append(cpes, cpeToString(*baseCPE))
		cpes = append(cpes, expandCPEsBySeparators(*baseCPE)...)
	}

	baseCPE.Vendor = baseCPE.Product
	cpes = append(cpes, cpeToString(*baseCPE))
	cpes = RemoveDuplicateCPES(cpes)

	// Retain base CPE
	if len(cpes) == 0 {
		cpes = append(cpes, cpeToString(*baseCPE))
	}

	return cpes
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
	}
	if strings.Contains(baseCPE.Vendor, "_") || strings.Contains(baseCPE.Product, "_") {
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
	}
	if strings.Contains(baseCPE.Vendor, ".") {
		tmp := baseCPE
		for _, vendor := range strings.Split(baseCPE.Vendor, ".") {
			baseCPE.Vendor = vendor
			cpes = append(cpes, cpeToString(baseCPE))
		}

		// Return to original value
		baseCPE = tmp
	}

	// Change Vendor to Product
	tmp := baseCPE
	baseCPE.Vendor = baseCPE.Product

	for _, vendor := range expand(baseCPE, "Vendor", '_', '-') {
		baseCPE.Vendor = vendor
		cpes = append(cpes, cpeToString(baseCPE))
		for _, product := range expand(baseCPE, "Product", '_', '-') {
			baseCPE.Product = product
			cpes = append(cpes, cpeToString(baseCPE))
		}
	}
	if strings.Contains(baseCPE.Vendor, "-") {
		baseCPE.Vendor = strings.Split(baseCPE.Vendor, "-")[0]
		for _, product := range expand(baseCPE, "Product", '_', '-') {
			baseCPE.Product = product
			cpes = append(cpes, cpeToString(baseCPE))
		}
	}
	if strings.Contains(baseCPE.Vendor, "_") {
		baseCPE.Vendor = strings.Split(baseCPE.Vendor, "_")[0]
		for _, product := range expand(baseCPE, "Product", '_', '-') {
			baseCPE.Product = product
			cpes = append(cpes, cpeToString(baseCPE))
		}
	}

	// Return to original value
	baseCPE = tmp

	return cpes
}

func expand(baseCPE CPE, f field, separator rune, replace rune) []string {
	expandedFields := make([]string, 0)
	switch f {
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
