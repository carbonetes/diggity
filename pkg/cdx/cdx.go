package cdx

import (
	"encoding/xml"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/CycloneDX/cyclonedx-go"
	diggity "github.com/carbonetes/diggity/internal/version"
)

var (
	// XMLN cyclonedx
	XMLN = fmt.Sprintf("http://cyclonedx.org/schema/bom/%+v", cyclonedx.SpecVersion1_5)
	lock sync.RWMutex
	BOM  *cyclonedx.BOM
)

const (
	cycloneDX = "CycloneDX"
	vendor    = "carbonetes"
	name      = "diggity"
)

func init() {
	BOM = New()
}

func New() *cyclonedx.BOM {
	return &cyclonedx.BOM{
		XMLName:     xml.Name{Local: cycloneDX},
		XMLNS:       XMLN,
		BOMFormat:   cycloneDX,
		Version:     1,
		SpecVersion: cyclonedx.SpecVersion1_5,
		Metadata:    getCDXMetadata(vendor, name, diggity.FromBuild().Version),
		Components:  &[]cyclonedx.Component{},
	}
}

func getCDXMetadata(author, name, version string) *cyclonedx.Metadata {
	return &cyclonedx.Metadata{
		Timestamp: time.Now().Format(time.RFC3339),
		Tools: &cyclonedx.ToolsChoice{
			Components: &[]cyclonedx.Component{
				{
					Type:    cyclonedx.ComponentTypeApplication,
					Author:  author,
					Name:    name,
					Version: version,
				},
			},
		},
	}
}

func AddComponent(c *cyclonedx.Component) {
	lock.Lock()
	defer lock.Unlock()

	if c == nil {
		return
	}

	// Check if the component already exists in the BOM
	for _, existingComponent := range *BOM.Components {
		if existingComponent.Name == c.Name && existingComponent.Version == c.Version {
			// If the component already exists, return without adding it
			return
		}
	}

	// If the component does not exist, add it to the BOM
	*BOM.Components = append(*BOM.Components, *c)
}

func SortComponents() {
	lock.Lock()
	defer lock.Unlock()

	// Sort components by name
	sort.Slice(*BOM.Components, func(i, j int) bool {
		return (*BOM.Components)[i].Name < (*BOM.Components)[j].Name
	})
}
