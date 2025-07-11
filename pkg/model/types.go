package model

import (
	"errors"
	"time"
)

// Component represents a software component in an SBOM
type Component struct {
	Name       string                 `json:"name"`
	Version    string                 `json:"version"`
	Type       string                 `json:"type"`
	PackageURL string                 `json:"purl,omitempty"`
	Licenses   []License              `json:"licenses,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// Package represents a detected software package
type Package struct {
	Name       string                 `json:"name"`
	Version    string                 `json:"version"`
	Type       string                 `json:"type"`
	PackageURL string                 `json:"purl,omitempty"`
	Language   string                 `json:"language,omitempty"`
	Licenses   []License              `json:"licenses,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// License represents a software license
type License struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	SPDXID    string    `json:"spdx_id,omitempty"`
	URL       string    `json:"url,omitempty"`
	Text      string    `json:"text,omitempty"`
	Component Component `json:"component"`
}

// Secret represents a detected secret or sensitive information
type Secret struct {
	Type        string    `json:"type"`
	Value       string    `json:"value,omitempty"` // May be redacted for security
	Description string    `json:"description"`
	Severity    Severity  `json:"severity"`
	Location    Location  `json:"location"`
	Component   Component `json:"component"`
}

// Location represents where a secret was found
type Location struct {
	File   string `json:"file,omitempty"`
	Line   int    `json:"line,omitempty"`
	Column int    `json:"column,omitempty"`
	Path   string `json:"path,omitempty"`
}

// Severity represents severity levels for issues (secrets, license violations, etc.)
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
	SeverityInfo     Severity = "info"
	SeverityUnknown  Severity = "unknown"
)

// AnalysisResult represents the result of SBOM component analysis (packages, licenses, secrets)
type AnalysisResult struct {
	ComponentsAnalyzed int                    `json:"components_analyzed"`
	Packages           []Package              `json:"packages"`
	Licenses           []License              `json:"licenses"`
	Secrets            []Secret               `json:"secrets"`
	Duration           time.Duration          `json:"duration"`
	Metadata           map[string]interface{} `json:"metadata"`
	Summary            *AnalysisSummary       `json:"summary,omitempty"`
}

// AnalysisSummary provides a high-level overview of analysis results
type AnalysisSummary struct {
	TotalPackages       int              `json:"total_packages"`
	TotalLicenses       int              `json:"total_licenses"`
	TotalSecrets        int              `json:"total_secrets"`
	SecretsBySeverity   map[Severity]int `json:"secrets_by_severity"`
	UniqueComponents    int              `json:"unique_components"`
	LicenseDistribution map[string]int   `json:"license_distribution"`
}

// ComplianceLevel represents overall compliance assessment levels
type ComplianceLevel string

const (
	ComplianceLevelCompliant    ComplianceLevel = "compliant"
	ComplianceLevelNonCompliant ComplianceLevel = "non_compliant"
	ComplianceLevelWarning      ComplianceLevel = "warning"
	ComplianceLevelUnknown      ComplianceLevel = "unknown"
)

// ComplianceAssessment represents comprehensive compliance evaluation results
type ComplianceAssessment struct {
	OverallCompliance ComplianceLevel        `json:"overall_compliance"`
	PolicyResults     []PolicyResult         `json:"policy_results"`
	Recommendations   []string               `json:"recommendations"`
	LicenseIssues     []LicenseIssue         `json:"license_issues,omitempty"`
	SecretIssues      []SecretIssue          `json:"secret_issues,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

// LicenseIssue represents a license compliance issue
type LicenseIssue struct {
	License     License   `json:"license"`
	Severity    Severity  `json:"severity"`
	Description string    `json:"description"`
	Component   Component `json:"component"`
}

// SecretIssue represents a secret detection issue
type SecretIssue struct {
	Secret      Secret    `json:"secret"`
	Severity    Severity  `json:"severity"`
	Description string    `json:"description"`
	Component   Component `json:"component"`
}

// PolicyResult represents the result of a policy evaluation
type PolicyResult struct {
	PolicyName string   `json:"policy_name"`
	Status     string   `json:"status"` // "pass", "fail", "warn"
	Message    string   `json:"message"`
	Violations []string `json:"violations,omitempty"`
}

// ScanResult represents the result of package scanning
type ScanResult struct {
	Packages        []Package              `json:"packages"`
	ComponentsFound int                    `json:"components_found"`
	Duration        time.Duration          `json:"duration"`
	Metadata        map[string]interface{} `json:"metadata"`
	Errors          []string               `json:"errors,omitempty"`
}

// SecretScanResult represents the result of secret scanning
type SecretScanResult struct {
	Secrets      []Secret               `json:"secrets"`
	SecretsFound int                    `json:"secrets_found"`
	Duration     time.Duration          `json:"duration"`
	Metadata     map[string]interface{} `json:"metadata"`
	Errors       []string               `json:"errors,omitempty"`
}

// LicenseScanResult represents the result of license scanning
type LicenseScanResult struct {
	Licenses      []License              `json:"licenses"`
	LicensesFound int                    `json:"licenses_found"`
	Duration      time.Duration          `json:"duration"`
	Metadata      map[string]interface{} `json:"metadata"`
	Errors        []string               `json:"errors,omitempty"`
}

// ComprehensiveScanResult represents the complete result of all scanning operations
type ComprehensiveScanResult struct {
	Target        interface{}            `json:"target"` // ScanTarget from scan package
	PackageResult *ScanResult            `json:"package_result,omitempty"`
	SecretResult  *SecretScanResult      `json:"secret_result,omitempty"`
	LicenseResult *LicenseScanResult     `json:"license_result,omitempty"`
	Duration      time.Duration          `json:"duration"`
	Timestamp     time.Time              `json:"timestamp"`
	Metadata      map[string]interface{} `json:"metadata"`
	Errors        []string               `json:"errors,omitempty"`
}

// Common errors
var (
	ErrUnsupportedSBOMFormat = errors.New("unsupported SBOM format")
	ErrInvalidComponent      = errors.New("invalid component data")
	ErrAnalysisTimeout       = errors.New("analysis timeout exceeded")
)
