package artifact

import (
	"github.com/carbonetes/diggity/pkg/model"
)

// Summary provides a high-level overview of the artifact graph
type Summary struct {
	TotalArtifacts      int                    `json:"total_artifacts"`
	TotalDependencies   int                    `json:"total_dependencies"`
	TotalLicenses       int                    `json:"total_licenses"`
	TotalSecrets        int                    `json:"total_secrets"`
	ArtifactsByType     map[string]int         `json:"artifacts_by_type"`
	LicenseDistribution map[string]int         `json:"license_distribution"`
	SecretsBySeverity   map[model.Severity]int `json:"secrets_by_severity"`
	RootArtifacts       int                    `json:"root_artifacts"`
	LeafArtifacts       int                    `json:"leaf_artifacts"`
}

// GetSummary returns a summary of the artifact graph
func (ag *ArtifactGraph) GetSummary() *Summary {
	summary := &Summary{
		TotalArtifacts:      len(ag.Artifacts),
		ArtifactsByType:     make(map[string]int),
		LicenseDistribution: make(map[string]int),
		SecretsBySeverity:   make(map[model.Severity]int),
	}

	// Count dependencies
	for _, deps := range ag.Dependencies {
		summary.TotalDependencies += len(deps)
	}

	// Count artifacts by type
	for _, artifact := range ag.Artifacts {
		summary.ArtifactsByType[artifact.Type]++
	}

	// Count licenses and their distribution
	allLicenses := ag.GetAllLicenses()
	summary.TotalLicenses = len(allLicenses)
	for _, license := range allLicenses {
		key := license.SPDXID
		if key == "" {
			key = license.Name
		}
		summary.LicenseDistribution[key]++
	}

	// Count secrets by severity
	allSecrets := ag.GetAllSecrets()
	summary.TotalSecrets = len(allSecrets)
	for _, secret := range allSecrets {
		summary.SecretsBySeverity[secret.Severity]++
	}

	// Count root and leaf artifacts
	summary.RootArtifacts = len(ag.GetRootArtifacts())
	summary.LeafArtifacts = len(ag.GetLeafArtifacts())

	return summary
}

// GetArtifactStats returns detailed statistics for a specific artifact
func (s *Summary) GetArtifactStats(artifactType string) int {
	return s.ArtifactsByType[artifactType]
}

// GetLicenseStats returns detailed statistics for a specific license
func (s *Summary) GetLicenseStats(licenseKey string) int {
	return s.LicenseDistribution[licenseKey]
}

// GetSecretStats returns detailed statistics for a specific severity level
func (s *Summary) GetSecretStats(severity model.Severity) int {
	return s.SecretsBySeverity[severity]
}

// HasCriticalSecrets returns true if there are any critical or high severity secrets
func (s *Summary) HasCriticalSecrets() bool {
	return s.SecretsBySeverity[model.SeverityCritical] > 0 || s.SecretsBySeverity[model.SeverityHigh] > 0
}

// GetMostCommonArtifactType returns the most common artifact type
func (s *Summary) GetMostCommonArtifactType() string {
	maxCount := 0
	mostCommon := ""
	for artifactType, count := range s.ArtifactsByType {
		if count > maxCount {
			maxCount = count
			mostCommon = artifactType
		}
	}
	return mostCommon
}

// GetMostCommonLicense returns the most common license
func (s *Summary) GetMostCommonLicense() string {
	maxCount := 0
	mostCommon := ""
	for license, count := range s.LicenseDistribution {
		if count > maxCount {
			maxCount = count
			mostCommon = license
		}
	}
	return mostCommon
}
