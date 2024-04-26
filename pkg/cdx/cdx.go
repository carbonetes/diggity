package cdx

import (
	"encoding/xml"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/CycloneDX/cyclonedx-go"
	diggity "github.com/carbonetes/diggity/internal/version"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/golistic/urn"
)

var (
	// XMLN cyclonedx
	XMLN = fmt.Sprintf("http://cyclonedx.org/schema/bom/%+v", cyclonedx.SpecVersion1_5)
	lock sync.RWMutex
	// BOM  *cyclonedx.BOM
)

const (
	cycloneDX = "CycloneDX"
	vendor    = "carbonetes"
	name      = "diggity"
)

// func init() {
// 	BOM = New()
// }

func New(addr *urn.URN) {
	stream.Set(addr.String(), &cyclonedx.BOM{
		XMLName:     xml.Name{Local: cycloneDX},
		XMLNS:       XMLN,
		BOMFormat:   cycloneDX,
		Version:     1,
		SpecVersion: cyclonedx.SpecVersion1_5,
		Metadata:    getCDXMetadata(vendor, name, diggity.FromBuild().Version),
		Components:  &[]cyclonedx.Component{},
	})
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

func AddComponent(c *cyclonedx.Component, addr *urn.URN) {
	if c == nil {
		return
	}

	data, _ := stream.Get(addr.String())
	bom, ok := data.(*cyclonedx.BOM)
	if !ok {
		log.Fatal("Failed to get BOM from stream")
	}

	// Check if the component already exists in the BOM
	for _, existingComponent := range *bom.Components {
		if existingComponent.Name == c.Name && existingComponent.Version == c.Version {
			// If the component already exists, return without adding it
			return
		}
	}

	// If the component does not exist, add it to the BOM
	// *BOM.Components = append(*BOM.Components, *c)
	*bom.Components = append(*bom.Components, *c)
	stream.Set(addr.String(), bom)
}

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
