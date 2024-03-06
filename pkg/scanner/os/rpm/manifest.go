package rpm

import (
	"strings"
)

// Format licenses
func formatLicenses(licensesGroup string) []string {
	if len(licensesGroup) > 0 && licensesGroup != " " {
		licenses := []string{}
		subgroups := strings.Split(licensesGroup, " and ")

		for _, group := range subgroups {
			group = strings.TrimSuffix(strings.TrimPrefix(group, "("), ")")
			licenses = append(licenses, strings.Split(group, " or ")...)
		}

		return deduplicateLicenses(licenses)
	}

	return []string{licensesGroup}
}

// Deduplicate licenses
func deduplicateLicenses(licenses []string) []string {
	uniqueLicenses := make(map[string]struct{})

	for _, license := range licenses {
		uniqueLicenses[license] = struct{}{}
	}

	deduplicated := make([]string, 0, len(uniqueLicenses))
	for license := range uniqueLicenses {
		deduplicated = append(deduplicated, license)
	}

	return deduplicated
}
