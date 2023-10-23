package maven

import (
	"regexp"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/model"
)

// Generate additional CPEs for java packages
func generateAdditionalCPE(vendor string, product string, version string, pkg *model.Package) {
	if len(vendor) > 0 {
		if strings.Contains(vendor, ".") {
			for _, v := range strings.Split(vendor, ".") {
				tldsRegex := `(com|org|io|edu|net|edu|gov|mil|the\ |a\ |an\ )(?:\b|')`
				if !regexp.MustCompile(tldsRegex).MatchString(v) {
					cpe.NewCPE23(pkg, v, product, version)
				}
			}
		} else {
			cpe.NewCPE23(pkg, vendor, product, version)
		}
	}
}
