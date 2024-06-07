package cpe

import (
	"github.com/CycloneDX/cyclonedx-go"
	"github.com/facebookincubator/nvdtools/wfn"
)

func Make(c *cyclonedx.Component) {
	if c.Type != cyclonedx.ComponentTypeLibrary {
		return
	}

	product := productLookup(c.Name)
	vendor := vendorLookup(product)
	if vendor == "" {
		vendor = product
	}

	cpe := new(product, vendor, c.Version)

	c.CPE = toString(cpe)

}

func new(product, vendor, version string) *wfn.Attributes {
	return &wfn.Attributes{
		Part:      "a",
		Vendor:    vendor,
		Product:   product,
		Version:   version,
		Update:    "*",
		Edition:   "*",
		Language:  "*",
		SWEdition: "*",
		TargetSW:  "*",
		TargetHW:  "*",
		Other:     "*",
	}
}

func join(matchers ...string) string {
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

func toString(cpe *wfn.Attributes) string {
	return join("cpe:2.3", cpe.Part, cpe.Vendor, cpe.Product, cpe.Version, cpe.Update, cpe.Edition, cpe.Language, cpe.SWEdition, cpe.TargetSW, cpe.TargetHW, cpe.Other)
}
