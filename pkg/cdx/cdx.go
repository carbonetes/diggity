package cdx

import (
	"encoding/xml"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/CycloneDX/cyclonedx-go"
	diggity "github.com/carbonetes/diggity/cmd/diggity/build"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx/dependency"
	"github.com/golistic/urn"
)

var (
	// XMLN cyclonedx namespace
	XMLN       = fmt.Sprintf("http://cyclonedx.org/schema/bom/%s", cyclonedx.SpecVersion1_5)
	mu         sync.RWMutex
	bomStorage = make(map[string]*cyclonedx.BOM)

	diggityVersion = diggity.FromBuild().Version
)

const (
	cycloneDX = "CycloneDX"
	vendor    = "carbonetes"
	author    = "Carbonetes Engineering Team"
	name      = "diggity"
	email     = "eng@carbonetes.com"
)

// New creates a new CycloneDX BOM
func New(addr *urn.URN) {
	mu.Lock()
	defer mu.Unlock()

	bomStorage[addr.String()] = &cyclonedx.BOM{
		XMLName:      xml.Name{Local: cycloneDX},
		XMLNS:        XMLN,
		BOMFormat:    cycloneDX,
		Version:      1,
		SpecVersion:  cyclonedx.SpecVersion1_5,
		Metadata:     setBasicMetadata(),
		Components:   &[]cyclonedx.Component{},
		Dependencies: &[]cyclonedx.Dependency{},
	}
}

// Clear removes the BOM from storage
func Clear(addr *urn.URN) {
	mu.Lock()
	defer mu.Unlock()
	delete(bomStorage, addr.String())
}

// AddComponent adds a single component to the BOM
func AddComponent(c *cyclonedx.Component, addr *urn.URN) {
	if c == nil {
		return
	}

	mu.Lock()
	defer mu.Unlock()

	bom := getBOMUnsafe(addr)
	*bom.Components = append(*bom.Components, *c)
}

// AddComponents adds multiple components to the BOM
func AddComponents(components *[]cyclonedx.Component, addr *urn.URN) {
	if components == nil || len(*components) == 0 {
		return
	}

	mu.Lock()
	defer mu.Unlock()

	bom := getBOMUnsafe(addr)
	*bom.Components = append(*bom.Components, *components...)
}

// SetMetadataComponent sets the metadata component
func SetMetadataComponent(addr *urn.URN, metadataComponent *cyclonedx.Component) {
	mu.Lock()
	defer mu.Unlock()

	bom := getBOMUnsafe(addr)
	bom.Metadata.Component = metadataComponent
}

// Finalize processes the BOM (deduplicate, sort, dependencies)
func Finalize(addr *urn.URN) *cyclonedx.BOM {
	mu.Lock()
	defer mu.Unlock()

	bom := getBOMUnsafe(addr)

	deduplicateComponents(bom)
	sortComponents(bom)
	parseDependencies(addr, bom)

	return bom
}

// getBOM safely retrieves BOM from storage
func getBOM(addr *urn.URN) *cyclonedx.BOM {
	mu.RLock()
	defer mu.RUnlock()
	return getBOMUnsafe(addr)
}

// getBOMUnsafe retrieves BOM without locking (for internal use when already locked)
func getBOMUnsafe(addr *urn.URN) *cyclonedx.BOM {
	bom, exists := bomStorage[addr.String()]
	if !exists {
		log.Fatal("BOM not found in storage for address: %s", addr.String())
	}
	return bom
}

// sortComponents sorts components by name
func sortComponents(bom *cyclonedx.BOM) {
	if bom.Components == nil || len(*bom.Components) == 0 {
		return
	}

	sort.Slice(*bom.Components, func(i, j int) bool {
		return (*bom.Components)[i].Name < (*bom.Components)[j].Name
	})
}

// deduplicateComponents removes duplicate components
func deduplicateComponents(bom *cyclonedx.BOM) {
	if bom.Components == nil || len(*bom.Components) == 0 {
		return
	}

	seen := make(map[string]bool)
	components := []cyclonedx.Component{}

	for _, c := range *bom.Components {
		key := fmt.Sprintf("%s-%s", c.Name, c.Version)
		if !seen[key] {
			components = append(components, c)
			seen[key] = true
		}
	}

	*bom.Components = components
}

// parseDependencies sets dependencies for components
func parseDependencies(addr *urn.URN, bom *cyclonedx.BOM) {
	dependencies := dependency.GetDependencyNodes(addr)
	if dependencies == nil {
		return
	}

	for i := range *dependencies {
		findDependencyRef(&(*dependencies)[i], bom.Components)
	}

	bom.Dependencies = dependencies
}

// findDependencyRef locates and replaces dependencies with BOMRefs
func findDependencyRef(node *cyclonedx.Dependency, components *[]cyclonedx.Component) {
	if node.Dependencies == nil || len(*node.Dependencies) == 0 {
		return
	}

	validDeps := []string{}

	for _, dep := range *node.Dependencies {
		for _, c := range *components {
			if c.Name == dep && c.BOMRef != "" {
				validDeps = append(validDeps, c.BOMRef)
				break
			}
		}
	}

	*node.Dependencies = validDeps
}

// setBasicMetadata creates basic metadata for the BOM
func setBasicMetadata() *cyclonedx.Metadata {
	return &cyclonedx.Metadata{
		Timestamp: time.Now().Format(time.RFC3339),
		Tools: &cyclonedx.ToolsChoice{
			Tools: &[]cyclonedx.Tool{
				{
					Vendor:  vendor,
					Name:    name,
					Version: diggityVersion,
				},
			},
		},
		Authors: &[]cyclonedx.OrganizationalContact{
			{
				Name:  author,
				Email: email,
			},
		},
	}
}
