package dependency

import (
	"sync"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/golistic/urn"
)

var lock *sync.RWMutex = &sync.RWMutex{}

const Type = "dependency"

func NewDependencyNodes(addr *urn.URN) {
	dependecyAddr := *addr
	dependecyAddr.NID = Type

	// Set the new map
	stream.Set(dependecyAddr.String(), &[]cyclonedx.Dependency{})
}

func AddDependency(addr *urn.URN, node *cyclonedx.Dependency) {
	lock.Lock()
	defer lock.Unlock()

	dependecyAddr := *addr
	dependecyAddr.NID = Type

	// Get the current map
	nodes := GetDependencyNodes(addr)
	if nodes == nil {
		log.Error("Dependency map not found")
		return
	}

	// Add the new node
	*nodes = append(*nodes, *node)

	// Set the new map
	stream.Set(dependecyAddr.String(), nodes)
}

func GetDependencyNodes(addr *urn.URN) *[]cyclonedx.Dependency {
	dependecyAddr := *addr
	dependecyAddr.NID = Type

	// Get the current map
	data, ok := stream.Get(dependecyAddr.String())
	if !ok {
		log.Error("Dependency map not found")
		return nil
	}

	nodes, ok := data.(*[]cyclonedx.Dependency)
	if !ok {
		log.Error("Dependency map is not a map")
		return nil
	}

	return nodes
}
