package model

import "github.com/carbonetes/jacked/pkg/core/model"

// PackageType represents the source/type of the package.
type PackageType string

const (
	OSPackage    PackageType = "os"
	ApplicationPackage PackageType = "application"
)

// Package actual package found
type Package struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Type            string                 `json:"type"`
	Version         string                 `json:"version"`
	PackageOrigin   PackageType            `json:"origin"`
	Parser          string                 `json:"parser"`
	Distro          string                 `json:"distro"`
	Language        string                 `json:"language"`
	Path            string                 `json:"path"`
	Locations       []Location             `json:"locations"`
	Description     string                 `json:"description,omitempty"`
	Licenses        []string               `json:"licenses,omitempty"`
	CPEs            []string               `json:"cpes"`
	PURL            PURL                   `json:"purl"`
	Metadata        interface{}            `json:"metadata"`
	Vulnerabilities *[]model.Vulnerability `json:"vulnerabilities,omitempty"`
}

// PURL - Package URL
type PURL string
