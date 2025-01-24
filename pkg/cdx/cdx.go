package cdx

import (
	"encoding/xml"
	"fmt"
	"sort"
	"sync"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/log"
	diggity "github.com/carbonetes/diggity/cmd/diggity/build"
	"github.com/carbonetes/diggity/pkg/cdx/dependency"
	stream "github.com/carbonetes/diggity/cmd/diggity/grove"
	"github.com/golistic/urn"
)

var (
	// XMLN cyclonedx
	XMLN = fmt.Sprintf("http://cyclonedx.org/schema/bom/%+v", cyclonedx.SpecVersion1_5)
	lock *sync.RWMutex

	diggityVersion = diggity.FromBuild().Version
)

const (
	cycloneDX = "CycloneDX"
	vendor    = "carbonetes"
	author    = "Carbonetes Engineering Team"
	name      = "diggity"
	email     = "eng@carbonetes.com"
)

func New(addr *urn.URN) {
	stream.Set(addr.String(), &cyclonedx.BOM{
		XMLName:      xml.Name{Local: cycloneDX},
		XMLNS:        XMLN,
		BOMFormat:    cycloneDX,
		Version:      1,
		SpecVersion:  cyclonedx.SpecVersion1_5,
		Metadata:     setBasicMetadata(),
		Components:   &[]cyclonedx.Component{},
		Dependencies: &[]cyclonedx.Dependency{},
	})
}

func AddComponent(c *cyclonedx.Component, addr *urn.URN) {
	if c == nil {
		return
	}

	data, _ := stream.Get(addr.String())
	bom, ok := data.(*cyclonedx.BOM)
	if !ok {
		log.Fatal("Failed to get BOM from stream")
	}

	*bom.Components = append(*bom.Components, *c)
	stream.Set(addr.String(), bom)
}

func AddComponents(components *[]cyclonedx.Component, addr *urn.URN) {
	if components == nil {
		return
	}

	data, _ := stream.Get(addr.String())
	bom, ok := data.(*cyclonedx.BOM)
	if !ok {
		log.Fatal("Failed to get BOM from stream")
	}

	*bom.Components = append(*bom.Components, *components...)
	stream.Set(addr.String(), bom)
}

func SetMetadataComponent(addr *urn.URN, metadataComponent *cyclonedx.Component) {
	data, _ := stream.Get(addr.String())
	bom := data.(*cyclonedx.BOM)

	bom.Metadata.Component = metadataComponent
	stream.Set(addr.String(), bom)
}

// Deprecated: Use Finalize() instead
func SortComponents(addr *urn.URN) *cyclonedx.BOM {
	lock.Lock()
	defer lock.Unlock()

	data, _ := stream.Get(addr.String())
	bom := data.(*cyclonedx.BOM)

	// Sort components by name
	sort.Slice(*bom.Components, func(i, j int) bool {
		return (*bom.Components)[i].Name < (*bom.Components)[j].Name
	})
	stream.Set(addr.String(), bom)
	return bom
}

func Finalize(addr *urn.URN) *cyclonedx.BOM {
	data, _ := stream.Get(addr.String())
	bom := data.(*cyclonedx.BOM)

	deduplicateComponents(bom)
	sortComponents(bom)
	parseDependencies(addr, bom)

	return bom
}

// Sort components by name
func sortComponents(bom *cyclonedx.BOM) {
	sort.Slice(*bom.Components, func(i, j int) bool {
		return (*bom.Components)[i].Name < (*bom.Components)[j].Name
	})
}

func deduplicateComponents(bom *cyclonedx.BOM) {
	seen := make(map[string]bool)
	components := []cyclonedx.Component{}
	for _, c := range *bom.Components {
		if _, ok := seen[c.Name]; !ok {
			components = append(components, c)
			seen[c.Name] = true
		}
	}
	*bom.Components = components
}

// Set Dependencies for each component in the BOM
func parseDependencies(addr *urn.URN, bom *cyclonedx.BOM) {
	dependencies := dependency.GetDependencyNodes(addr)
	if dependencies != nil {
		for _, d := range *dependencies {
			findDependencyRef(&d, bom.Components)
		}
	}
	bom.Dependencies = dependencies
}

// Locate and replace dependencies with BOMRefs
func findDependencyRef(node *cyclonedx.Dependency, components *[]cyclonedx.Component) {
	toBeRemoved := []int{}
	for i, dep := range *node.Dependencies {
		found := new(string)
		for _, c := range *components {
			if c.Name == dep {
				found = &c.BOMRef
				break
			}
		}
		if *found != "" {
			(*node.Dependencies)[i] = *found
		} else {
			toBeRemoved = append(toBeRemoved, i)
		}
	}

	// Remove dependencies that are not found in the components
	for i := len(toBeRemoved) - 1; i >= 0; i-- {
		(*node.Dependencies) = append((*node.Dependencies)[:toBeRemoved[i]], (*node.Dependencies)[toBeRemoved[i]+1:]...)
	}
}
