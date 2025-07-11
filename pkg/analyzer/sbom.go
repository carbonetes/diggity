package analyzer

import (
	"context"
	"time"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/config"
	"github.com/carbonetes/diggity/pkg/model"
)

// Analyzer interface for SBOM-specific operations (following terminology refactoring)
type Analyzer interface {
	// AnalyzeSBOM performs component analysis on structured SBOMs (packages, licenses, secrets)
	AnalyzeSBOM(ctx context.Context, sbom interface{}) (*model.AnalysisResult, error)

	// AnalyzeComponents processes individual components for packages, licenses, and secrets
	AnalyzeComponents(ctx context.Context, components []model.Component) (*model.AnalysisResult, error)

	// AnalyzeLicenseCompliance evaluates license compliance and policy adherence
	AnalyzeLicenseCompliance(ctx context.Context, result *model.AnalysisResult, policies *config.PolicyConfig) (*model.ComplianceAssessment, error)
}

// SBOMAnalyzer implements the Analyzer interface for CycloneDX SBOMs
type SBOMAnalyzer struct {
	config *config.AnalysisConfig
}

// NewSBOMAnalyzer creates a new SBOM analyzer instance
func NewSBOMAnalyzer(cfg *config.AnalysisConfig) *SBOMAnalyzer {
	return &SBOMAnalyzer{
		config: cfg,
	}
}

// AnalyzeSBOM performs comprehensive component analysis on a CycloneDX SBOM
// This follows the terminology refactoring: "Analyzing" for structured SBOM processing
// Focus: packages, licenses, and secrets detection
func (a *SBOMAnalyzer) AnalyzeSBOM(ctx context.Context, sbom interface{}) (*model.AnalysisResult, error) {
	log.Debug("Starting SBOM component analysis (packages, licenses, secrets)")
	start := time.Now()

	// Type assertion for CycloneDX BOM
	bom, ok := sbom.(*cyclonedx.BOM)
	if !ok {
		return nil, model.ErrUnsupportedSBOMFormat
	}

	if bom == nil {
		log.Debug("BOM is nil, skipping analysis")
		return &model.AnalysisResult{}, nil
	}

	if bom.Components == nil || len(*bom.Components) == 0 {
		log.Debug("BOM has no components, skipping analysis")
		return &model.AnalysisResult{
			ComponentsAnalyzed: 0,
			Duration:           time.Since(start),
		}, nil
	}

	components := a.convertCycloneDXComponents(*bom.Components)
	log.Debugf("Analyzing %d components for packages, licenses, and secrets", len(components))

	result, err := a.AnalyzeComponents(ctx, components)
	if err != nil {
		return nil, err
	}

	result.Duration = time.Since(start)
	log.Debugf("SBOM component analysis completed in %v", result.Duration)

	return result, nil
}

// AnalyzeComponents processes individual components for packages, licenses, and secrets
func (a *SBOMAnalyzer) AnalyzeComponents(ctx context.Context, components []model.Component) (*model.AnalysisResult, error) {
	result := &model.AnalysisResult{
		ComponentsAnalyzed: len(components),
		Packages:           make([]model.Package, 0),
		Licenses:           make([]model.License, 0),
		Secrets:            make([]model.Secret, 0),
		Metadata:           make(map[string]interface{}),
	}

	// Process each component for package information, license detection, and secret scanning
	for _, component := range components {
		// Extract package information
		pkg := model.Package{
			Name:       component.Name,
			Version:    component.Version,
			Type:       component.Type,
			PackageURL: component.PackageURL,
		}
		result.Packages = append(result.Packages, pkg)

		// Analyze licenses (this would integrate with license detection)
		if len(component.Licenses) > 0 {
			for _, license := range component.Licenses {
				result.Licenses = append(result.Licenses, license)
			}
		}

		// Scan for secrets (this would integrate with secret detection)
		secrets := a.scanComponentForSecrets(component)
		result.Secrets = append(result.Secrets, secrets...)
	}

	return result, nil
}

// scanComponentForSecrets scans a component for potential secrets
func (a *SBOMAnalyzer) scanComponentForSecrets(component model.Component) []model.Secret {
	var secrets []model.Secret

	// Scan component metadata for potential secrets
	for key, value := range component.Metadata {
		if valueStr, ok := value.(string); ok {
			if a.containsPotentialSecret(key, valueStr) {
				secret := model.Secret{
					Type:        "metadata_secret",
					Description: "Potential secret found in component metadata",
					Severity:    a.assessSecretSeverity(key, valueStr),
					Location: model.Location{
						Path: component.Name,
					},
					Component: component,
				}
				secrets = append(secrets, secret)
			}
		}
	}

	return secrets
}

// containsPotentialSecret checks if a key-value pair contains potential secrets
func (a *SBOMAnalyzer) containsPotentialSecret(key, value string) bool {
	if a.config == nil {
		return false
	}

	// Check against configured secret patterns
	for _, pattern := range a.config.SecretPatterns {
		// Simple pattern matching - in real implementation would use regex
		if len(value) > 20 && (key == pattern.Name || key == "token" || key == "key" || key == "password" || key == "secret") {
			return true
		}
	}

	return false
}

// assessSecretSeverity determines the severity of a potential secret
func (a *SBOMAnalyzer) assessSecretSeverity(key, value string) model.Severity {
	// Determine severity based on key type and value characteristics
	switch key {
	case "password", "private_key":
		return model.SeverityCritical
	case "token", "api_key":
		return model.SeverityHigh
	case "secret":
		return model.SeverityMedium
	default:
		return model.SeverityLow
	}
}

// AnalyzeLicenseCompliance evaluates license compliance and policy adherence
func (a *SBOMAnalyzer) AnalyzeLicenseCompliance(ctx context.Context, result *model.AnalysisResult, policies *config.PolicyConfig) (*model.ComplianceAssessment, error) {
	assessment := &model.ComplianceAssessment{
		OverallCompliance: model.ComplianceLevelCompliant,
		PolicyResults:     make([]model.PolicyResult, 0),
		Recommendations:   make([]string, 0),
		LicenseIssues:     make([]model.LicenseIssue, 0),
		SecretIssues:      make([]model.SecretIssue, 0),
	}

	// Evaluate license compliance
	for _, license := range result.Licenses {
		if a.isLicenseProblematic(license, policies) {
			issue := model.LicenseIssue{
				License:     license,
				Severity:    model.SeverityMedium,
				Description: "License may not comply with organizational policies",
				Component:   license.Component,
			}
			assessment.LicenseIssues = append(assessment.LicenseIssues, issue)
		}
	}

	// Evaluate secret issues
	for _, secret := range result.Secrets {
		if secret.Severity == model.SeverityCritical || secret.Severity == model.SeverityHigh {
			issue := model.SecretIssue{
				Secret:      secret,
				Severity:    secret.Severity,
				Description: "Potential secret detected in component",
				Component:   secret.Component,
			}
			assessment.SecretIssues = append(assessment.SecretIssues, issue)
		}
	}

	// Determine overall compliance based on issues found
	if len(assessment.LicenseIssues) > 0 || len(assessment.SecretIssues) > 0 {
		assessment.OverallCompliance = model.ComplianceLevelWarning
		if a.hasCriticalIssues(assessment) {
			assessment.OverallCompliance = model.ComplianceLevelNonCompliant
		}
	}

	return assessment, nil
}

// isLicenseProblematic checks if a license violates policies
func (a *SBOMAnalyzer) isLicenseProblematic(license model.License, policies *config.PolicyConfig) bool {
	if policies == nil {
		return false
	}

	// Check against license policies
	for _, policy := range policies.LicensePolicies {
		if a.isLicenseForbidden(license, policy) {
			return true
		}

		if a.isLicenseNotAllowed(license, policy) {
			return true
		}
	}

	return false
}

// isLicenseForbidden checks if license is in forbidden list
func (a *SBOMAnalyzer) isLicenseForbidden(license model.License, policy config.LicensePolicy) bool {
	for _, forbidden := range policy.ForbiddenSPDX {
		if license.SPDXID == forbidden || license.ID == forbidden {
			return true
		}
	}
	return false
}

// isLicenseNotAllowed checks if license is not in allowed list (when specified)
func (a *SBOMAnalyzer) isLicenseNotAllowed(license model.License, policy config.LicensePolicy) bool {
	if len(policy.AllowedSPDX) == 0 {
		return false
	}

	for _, allowed := range policy.AllowedSPDX {
		if license.SPDXID == allowed || license.ID == allowed {
			return false
		}
	}
	return true
}

// hasCriticalIssues determines if there are critical compliance issues
func (a *SBOMAnalyzer) hasCriticalIssues(assessment *model.ComplianceAssessment) bool {
	for _, issue := range assessment.SecretIssues {
		if issue.Severity == model.SeverityCritical {
			return true
		}
	}
	return false
}

// convertCycloneDXComponents converts CycloneDX components to internal model
func (a *SBOMAnalyzer) convertCycloneDXComponents(cdxComponents []cyclonedx.Component) []model.Component {
	components := make([]model.Component, 0, len(cdxComponents))

	for _, cdxComp := range cdxComponents {
		comp := model.Component{
			Name:    cdxComp.Name,
			Version: cdxComp.Version,
			Type:    string(cdxComp.Type),
		}

		if cdxComp.PackageURL != "" {
			comp.PackageURL = cdxComp.PackageURL
		}

		components = append(components, comp)
	}

	return components
}
