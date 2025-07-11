package config

import "time"

// AnalysisConfig contains configuration for SBOM analysis
type AnalysisConfig struct {
	// Analysis timeout settings
	AnalysisTimeout time.Duration `yaml:"analysis_timeout" json:"analysis_timeout"`

	// Concurrent analyzers (following terminology refactoring)
	MaxConcurrentAnalyzers int `yaml:"max_concurrent_analyzers" json:"max_concurrent_analyzers"`

	// Progressive analysis (following terminology refactoring)
	EnableProgressiveAnalysis bool `yaml:"enable_progressive_analysis" json:"enable_progressive_analysis"`

	// Secret scanning configuration
	SecretScanningEnabled bool            `yaml:"secret_scanning_enabled" json:"secret_scanning_enabled"`
	SecretPatterns        []SecretPattern `yaml:"secret_patterns" json:"secret_patterns"`

	// License analysis configuration
	LicenseAnalysisEnabled bool     `yaml:"license_analysis_enabled" json:"license_analysis_enabled"`
	AllowedLicenses        []string `yaml:"allowed_licenses" json:"allowed_licenses"`
	ForbiddenLicenses      []string `yaml:"forbidden_licenses" json:"forbidden_licenses"`

	// Package analysis configuration
	PackageAnalysisEnabled bool     `yaml:"package_analysis_enabled" json:"package_analysis_enabled"`
	IncludeDevDependencies bool     `yaml:"include_dev_dependencies" json:"include_dev_dependencies"`
	ExcludePatterns        []string `yaml:"exclude_patterns" json:"exclude_patterns"`
}

// SecretPattern defines patterns for secret detection
type SecretPattern struct {
	Name        string `yaml:"name" json:"name"`
	Pattern     string `yaml:"pattern" json:"pattern"`
	Description string `yaml:"description" json:"description"`
	Severity    string `yaml:"severity" json:"severity"`
}

// PolicyConfig contains policy configuration for compliance assessment
type PolicyConfig struct {
	// License policies
	LicensePolicies []LicensePolicy `yaml:"license_policies" json:"license_policies"`

	// Secret policies
	SecretPolicies []SecretPolicy `yaml:"secret_policies" json:"secret_policies"`

	// General policies
	EnforceCompliance bool `yaml:"enforce_compliance" json:"enforce_compliance"`
	FailOnViolation   bool `yaml:"fail_on_violation" json:"fail_on_violation"`
}

// LicensePolicy defines license compliance rules
type LicensePolicy struct {
	Name          string   `yaml:"name" json:"name"`
	AllowedSPDX   []string `yaml:"allowed_spdx" json:"allowed_spdx"`
	ForbiddenSPDX []string `yaml:"forbidden_spdx" json:"forbidden_spdx"`
	Severity      string   `yaml:"severity" json:"severity"`
	Description   string   `yaml:"description" json:"description"`
}

// SecretPolicy defines secret detection and handling rules
type SecretPolicy struct {
	Name        string   `yaml:"name" json:"name"`
	Patterns    []string `yaml:"patterns" json:"patterns"`
	Severity    string   `yaml:"severity" json:"severity"`
	Action      string   `yaml:"action" json:"action"` // "warn", "fail", "ignore"
	Description string   `yaml:"description" json:"description"`
}

// ScanConfig contains configuration for raw scanning operations
type ScanConfig struct {
	// Scanning timeout settings
	ScanTimeout time.Duration `yaml:"scan_timeout" json:"scan_timeout"`

	// Concurrent scanners (following terminology refactoring)
	MaxConcurrentScanners int `yaml:"max_concurrent_scanners" json:"max_concurrent_scanners"`

	// Progressive scanning (following terminology refactoring)
	EnableProgressiveScanning bool `yaml:"enable_progressive_scanning" json:"enable_progressive_scanning"`

	// Package scanning configuration
	PackageScanningEnabled bool     `yaml:"package_scanning_enabled" json:"package_scanning_enabled"`
	SupportedPackageTypes  []string `yaml:"supported_package_types" json:"supported_package_types"`

	// Secret scanning configuration
	SecretScanningEnabled bool            `yaml:"secret_scanning_enabled" json:"secret_scanning_enabled"`
	SecretPatterns        []SecretPattern `yaml:"secret_patterns" json:"secret_patterns"`

	// License scanning configuration
	LicenseScanningEnabled bool `yaml:"license_scanning_enabled" json:"license_scanning_enabled"`
	LicenseDetectionDepth  int  `yaml:"license_detection_depth" json:"license_detection_depth"`

	// Performance settings
	MaxFileSize     int64    `yaml:"max_file_size" json:"max_file_size"`
	ExcludePatterns []string `yaml:"exclude_patterns" json:"exclude_patterns"`
	IncludePatterns []string `yaml:"include_patterns" json:"include_patterns"`
}

// DefaultAnalysisConfig returns default configuration for analysis
func DefaultAnalysisConfig() *AnalysisConfig {
	return &AnalysisConfig{
		AnalysisTimeout:           5 * time.Minute,
		MaxConcurrentAnalyzers:    4,
		EnableProgressiveAnalysis: true,
		SecretScanningEnabled:     true,
		LicenseAnalysisEnabled:    true,
		PackageAnalysisEnabled:    true,
		IncludeDevDependencies:    false,
		SecretPatterns: []SecretPattern{
			{
				Name:        "API Key",
				Pattern:     `(?i)api[_-]?key.*[=:]\s*['""]?([a-zA-Z0-9]{20,})['""]?`,
				Description: "Generic API key pattern",
				Severity:    "high",
			},
			{
				Name:        "Private Key",
				Pattern:     `-----BEGIN\s+(?:RSA\s+)?PRIVATE\s+KEY-----`,
				Description: "Private key detection",
				Severity:    "critical",
			},
		},
		AllowedLicenses: []string{"MIT", "Apache-2.0", "BSD-3-Clause"},
		ExcludePatterns: []string{"**/test/**", "**/tests/**", "**/*_test.go"},
	}
}

// DefaultPolicyConfig returns default policy configuration
func DefaultPolicyConfig() *PolicyConfig {
	return &PolicyConfig{
		EnforceCompliance: true,
		FailOnViolation:   false,
		LicensePolicies: []LicensePolicy{
			{
				Name:          "Open Source Compliant",
				AllowedSPDX:   []string{"MIT", "Apache-2.0", "BSD-3-Clause", "ISC"},
				ForbiddenSPDX: []string{"GPL-3.0", "AGPL-3.0"},
				Severity:      "medium",
				Description:   "Ensures only approved open source licenses are used",
			},
		},
		SecretPolicies: []SecretPolicy{
			{
				Name:        "No Hardcoded Secrets",
				Patterns:    []string{"password", "secret", "token", "key"},
				Severity:    "high",
				Action:      "fail",
				Description: "Prevents hardcoded secrets in components",
			},
		},
	}
}

// DefaultScanConfig returns default configuration for scanning
func DefaultScanConfig() *ScanConfig {
	return &ScanConfig{
		ScanTimeout:               10 * time.Minute,
		MaxConcurrentScanners:     4,
		EnableProgressiveScanning: true,
		PackageScanningEnabled:    true,
		SecretScanningEnabled:     true,
		LicenseScanningEnabled:    true,
		LicenseDetectionDepth:     3,
		MaxFileSize:               100 * 1024 * 1024, // 100MB
		SupportedPackageTypes: []string{
			"npm", "pypi", "maven", "gradle", "cargo", "go", "composer",
			"rubygem", "nuget", "swift", "cocoapods", "hex", "pub", "conan",
		},
		SecretPatterns: []SecretPattern{
			{
				Name:        "API Key",
				Pattern:     `(?i)api[_-]?key.*[=:]\s*['""]?([a-zA-Z0-9]{20,})['""]?`,
				Description: "Generic API key pattern",
				Severity:    "high",
			},
		},
		ExcludePatterns: []string{
			"**/test/**", "**/tests/**", "**/*_test.go", "**/node_modules/**",
			"**/vendor/**", "**/.git/**", "**/target/**", "**/build/**",
		},
	}
}
