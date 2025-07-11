package artifact

import (
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
)

func TestArtifactCreation(t *testing.T) {
	// Create a test package
	pkg := &model.Package{
		Name:    "test-package",
		Version: "1.0.0",
		Type:    "npm",
	}

	// Create artifact
	artifact := NewArtifact("test-artifact-1", pkg)

	if artifact.ID != "test-artifact-1" {
		t.Errorf("Expected artifact ID 'test-artifact-1', got '%s'", artifact.ID)
	}

	if artifact.Package.Name != "test-package" {
		t.Errorf("Expected package name 'test-package', got '%s'", artifact.Package.Name)
	}
}

func TestDependencyCreation(t *testing.T) {
	dep := NewDependency("lodash", "4.17.21", DependencyTypeDirect)

	if dep.Name != "lodash" {
		t.Errorf("Expected dependency name 'lodash', got '%s'", dep.Name)
	}

	if dep.Version != "4.17.21" {
		t.Errorf("Expected dependency version '4.17.21', got '%s'", dep.Version)
	}

	if !dep.IsDirectDependency() {
		t.Error("Expected dependency to be direct")
	}
}

func TestArtifactGraph(t *testing.T) {
	graph := NewArtifactGraph()

	// Create test artifacts
	pkg1 := &model.Package{Name: "pkg1", Version: "1.0.0", Type: "npm"}
	pkg2 := &model.Package{Name: "pkg2", Version: "2.0.0", Type: "npm"}

	artifact1 := NewArtifact("artifact1", pkg1)
	artifact2 := NewArtifact("artifact2", pkg2)

	// Add artifacts to graph
	graph.AddArtifact(artifact1)
	graph.AddArtifact(artifact2)

	// Add dependency relationship
	graph.AddDependency("artifact1", "artifact2")

	// Test graph functionality
	if len(graph.GetAllArtifacts()) != 2 {
		t.Errorf("Expected 2 artifacts in graph, got %d", len(graph.GetAllArtifacts()))
	}

	deps := graph.GetDependencies("artifact1")
	if len(deps) != 1 || deps[0] != "artifact2" {
		t.Errorf("Expected artifact1 to depend on artifact2")
	}

	dependents := graph.GetDependents("artifact2")
	if len(dependents) != 1 || dependents[0] != "artifact1" {
		t.Errorf("Expected artifact2 to have artifact1 as dependent")
	}
}

func TestArtifactSummary(t *testing.T) {
	graph := NewArtifactGraph()

	// Create test artifact with license and secret
	pkg := &model.Package{Name: "test-pkg", Version: "1.0.0", Type: "npm"}
	artifact := NewArtifact("test-artifact", pkg)
	artifact.Type = "npm"

	// Add license
	license := model.License{
		Name:   "MIT",
		SPDXID: "MIT",
	}
	artifact.AddLicense(license)

	// Add secret
	secret := model.Secret{
		Type:     "api-key",
		Value:    "secret-value",
		Severity: model.SeverityHigh,
	}
	artifact.AddSecret(secret)

	graph.AddArtifact(artifact)

	// Get summary
	summary := graph.GetSummary()

	if summary.TotalArtifacts != 1 {
		t.Errorf("Expected 1 artifact, got %d", summary.TotalArtifacts)
	}

	if summary.TotalLicenses != 1 {
		t.Errorf("Expected 1 license, got %d", summary.TotalLicenses)
	}

	if summary.TotalSecrets != 1 {
		t.Errorf("Expected 1 secret, got %d", summary.TotalSecrets)
	}

	if !summary.HasCriticalSecrets() {
		t.Error("Expected to have critical secrets")
	}
}
