package artifact

import (
	"github.com/carbonetes/diggity/pkg/model"
)

// ArtifactGraph manages artifact dependency mapping and relationship tracking
type ArtifactGraph struct {
	Artifacts    map[string]*Artifact       `json:"artifacts"`    // artifact ID -> artifact
	Dependencies map[string][]string        `json:"dependencies"` // artifact ID -> dependency IDs
	Dependents   map[string][]string        `json:"dependents"`   // artifact ID -> dependent IDs (reverse mapping)
	Licenses     map[string][]model.License `json:"licenses"`     // license tracking across artifacts
	Secrets      map[string][]model.Secret  `json:"secrets"`      // secret tracking across artifacts
	Metadata     map[string]interface{}     `json:"metadata"`     // Graph-level metadata
}

// NewArtifactGraph creates a new artifact dependency graph
func NewArtifactGraph() *ArtifactGraph {
	return &ArtifactGraph{
		Artifacts:    make(map[string]*Artifact),
		Dependencies: make(map[string][]string),
		Dependents:   make(map[string][]string),
		Licenses:     make(map[string][]model.License),
		Secrets:      make(map[string][]model.Secret),
		Metadata:     make(map[string]interface{}),
	}
}

// AddArtifact adds an artifact to the graph and tracks its licenses and secrets
func (ag *ArtifactGraph) AddArtifact(artifact *Artifact) {
	ag.Artifacts[artifact.ID] = artifact

	// Track licenses for this artifact
	if len(artifact.Licenses) > 0 {
		ag.Licenses[artifact.ID] = artifact.Licenses
	}

	// Track secrets for this artifact
	if len(artifact.Secrets) > 0 {
		ag.Secrets[artifact.ID] = artifact.Secrets
	}
}

// AddDependency adds a dependency relationship between artifacts
func (ag *ArtifactGraph) AddDependency(artifactID, dependencyID string) {
	// Add to dependencies mapping
	if ag.Dependencies[artifactID] == nil {
		ag.Dependencies[artifactID] = make([]string, 0)
	}
	ag.Dependencies[artifactID] = append(ag.Dependencies[artifactID], dependencyID)

	// Add to dependents mapping (reverse relationship)
	if ag.Dependents[dependencyID] == nil {
		ag.Dependents[dependencyID] = make([]string, 0)
	}
	ag.Dependents[dependencyID] = append(ag.Dependents[dependencyID], artifactID)
}

// GetDependencies returns all dependencies for an artifact
func (ag *ArtifactGraph) GetDependencies(artifactID string) []string {
	return ag.Dependencies[artifactID]
}

// GetDependents returns all dependents for an artifact
func (ag *ArtifactGraph) GetDependents(artifactID string) []string {
	return ag.Dependents[artifactID]
}

// GetArtifact returns an artifact by ID
func (ag *ArtifactGraph) GetArtifact(artifactID string) *Artifact {
	return ag.Artifacts[artifactID]
}

// GetAllArtifacts returns all artifacts in the graph
func (ag *ArtifactGraph) GetAllArtifacts() []*Artifact {
	artifacts := make([]*Artifact, 0, len(ag.Artifacts))
	for _, artifact := range ag.Artifacts {
		artifacts = append(artifacts, artifact)
	}
	return artifacts
}

// GetAllLicenses returns all unique licenses across artifacts
func (ag *ArtifactGraph) GetAllLicenses() []model.License {
	licenseMap := make(map[string]model.License)
	for _, licenses := range ag.Licenses {
		for _, license := range licenses {
			// Use SPDX ID as key, fallback to Name if SPDX ID is empty
			key := license.SPDXID
			if key == "" {
				key = license.Name
			}
			licenseMap[key] = license
		}
	}

	result := make([]model.License, 0, len(licenseMap))
	for _, license := range licenseMap {
		result = append(result, license)
	}
	return result
}

// GetAllSecrets returns all secrets across artifacts
func (ag *ArtifactGraph) GetAllSecrets() []model.Secret {
	var allSecrets []model.Secret
	for _, secrets := range ag.Secrets {
		allSecrets = append(allSecrets, secrets...)
	}
	return allSecrets
}

// GetLicensesByArtifact returns licenses for a specific artifact
func (ag *ArtifactGraph) GetLicensesByArtifact(artifactID string) []model.License {
	return ag.Licenses[artifactID]
}

// GetSecretsByArtifact returns secrets for a specific artifact
func (ag *ArtifactGraph) GetSecretsByArtifact(artifactID string) []model.Secret {
	return ag.Secrets[artifactID]
}

// GetDependencyTree returns the full dependency tree starting from an artifact
func (ag *ArtifactGraph) GetDependencyTree(artifactID string, visited map[string]bool) map[string]*Artifact {
	if visited == nil {
		visited = make(map[string]bool)
	}

	tree := make(map[string]*Artifact)

	// Avoid circular dependencies
	if visited[artifactID] {
		return tree
	}
	visited[artifactID] = true

	// Add current artifact
	if artifact := ag.GetArtifact(artifactID); artifact != nil {
		tree[artifactID] = artifact
	}

	// Recursively add dependencies
	for _, depID := range ag.GetDependencies(artifactID) {
		depTree := ag.GetDependencyTree(depID, visited)
		for id, artifact := range depTree {
			tree[id] = artifact
		}
	}

	return tree
}

// GetRootArtifacts returns artifacts that have no dependents (root level packages)
func (ag *ArtifactGraph) GetRootArtifacts() []*Artifact {
	var roots []*Artifact
	for id, artifact := range ag.Artifacts {
		if len(ag.Dependents[id]) == 0 {
			roots = append(roots, artifact)
		}
	}
	return roots
}

// GetLeafArtifacts returns artifacts that have no dependencies (leaf level packages)
func (ag *ArtifactGraph) GetLeafArtifacts() []*Artifact {
	var leaves []*Artifact
	for id, artifact := range ag.Artifacts {
		if len(ag.Dependencies[id]) == 0 {
			leaves = append(leaves, artifact)
		}
	}
	return leaves
}
