package rpm

import (
	"fmt"
	"strings"

	"github.com/carbonetes/diggity/pkg/types"
	"github.com/google/uuid"
	rpmdb "github.com/knqyf263/go-rpmdb/pkg"
)

// Produce a component from a package info
func newComponent(info rpmdb.PackageInfo) types.Component{
	return types.Component{
		ID:      uuid.NewString(),
		Name:    info.Name,
		Type:    Type,
		Version: fmt.Sprintf("%+v-%+v", info.Version, info.Release),
		Description: info.Summary,
		Licenses: formatLicenses(info.License),
		PURL: fmt.Sprintf("pkg:%s/%s@%s?arch=%s", Type, info.Name, info.Version, info.Arch),
		Metadata: info,
	}
}

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
