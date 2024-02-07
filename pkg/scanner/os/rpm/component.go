package rpm

import (
	"fmt"
	"strings"

	"github.com/carbonetes/diggity/pkg/types"
	"github.com/google/uuid"
	rpmdb "github.com/knqyf263/go-rpmdb/pkg"
)

type Metadata struct {
	Epoch           *int   `json:"epoch,omitempty"`
	Name            string `json:"name,omitempty"`
	Version         string `json:"version,omitempty"`
	Release         string `json:"release,omitempty"`
	Arch            string `json:"arch,omitempty"`
	SourceRpm       string `json:"sourceRpm,omitempty"`
	Size            int    `json:"size,omitempty"`
	License         string `json:"license,omitempty"`
	Vendor          string `json:"vendor,omitempty"`
	Modularitylabel string `json:"modularitylabel,omitempty"`
	Summary         string `json:"summary,omitempty"`
	PGP             string `json:"pgp,omitempty"`
	SigMD5          string `json:"sigMD5,omitempty"`
	DigestAlgorithm int    `json:"digestAlgorithm,omitempty"`
	InstallTime     int    `json:"installTime,omitempty"`
}

// Produce a component from a package info
func newComponent(info rpmdb.PackageInfo) types.Component {

	metadata := Metadata{
		Epoch:           info.Epoch,
		Name:            info.Name,
		Version:         info.Version,
		Release:         info.Release,
		Arch:            info.Arch,
		SourceRpm:       info.SourceRpm,
		Size:            info.Size,
		License:         info.License,
		Vendor:          info.Vendor,
		Modularitylabel: info.Modularitylabel,
		Summary:         info.Summary,
		PGP:             info.PGP,
		SigMD5:          info.SigMD5,
		DigestAlgorithm: int(info.DigestAlgorithm),
		InstallTime:     info.InstallTime,
	}

	return types.Component{
		ID:          uuid.NewString(),
		Name:        info.Name,
		Type:        Type,
		Version:     fmt.Sprintf("%+v-%+v", info.Version, info.Release),
		Description: info.Summary,
		Licenses:    formatLicenses(info.License),
		PURL:        fmt.Sprintf("pkg:%s/%s@%s?arch=%s", Type, info.Name, info.Version, info.Arch),
		Metadata:    metadata,
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
