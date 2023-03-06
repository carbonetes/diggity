package cpe

import (
	"fmt"
	"regexp"

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
