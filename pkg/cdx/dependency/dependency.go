package dependency

import (
	"sync"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/golistic/urn"
)

var lock *sync.RWMutex = &sync.RWMutex{}
var storage = make(map[string]*[]cyclonedx.Dependency)

const Type = "dependency"

func NewDependencyNodes(addr *urn.URN) {
	lock.Lock()
	defer lock.Unlock()

	dependecyAddr := *addr
	dependecyAddr.NID = Type

	// Set the new map
	storage[dependecyAddr.String()] = &[]cyclonedx.Dependency{}
}

func AddDependency(addr *urn.URN, node *cyclonedx.Dependency) {
	lock.Lock()
	defer lock.Unlock()

	dependecyAddr := *addr
	dependecyAddr.NID = Type

	// Get the current map
	nodes := getDependencyNodesUnsafe(&dependecyAddr)
	if nodes == nil {
		log.Debug("Dependency map not found")
		return
	}

	// Add the new node
	*nodes = append(*nodes, *node)

	// Store the updated map
	storage[dependecyAddr.String()] = nodes
}

func GetDependencyNodes(addr *urn.URN) *[]cyclonedx.Dependency {
	lock.RLock()
	defer lock.RUnlock()

	dependecyAddr := *addr
	dependecyAddr.NID = Type

	return getDependencyNodesUnsafe(&dependecyAddr)
}

func getDependencyNodesUnsafe(addr *urn.URN) *[]cyclonedx.Dependency {
	// Get the current map
	nodes, ok := storage[addr.String()]
	if !ok {
		log.Debug("Dependency map not found")
		return nil
	}

	return nodes
}
