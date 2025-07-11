package registry

import (
	"fmt"
	"sync"

	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/scan/types"
)

// Registry manages all available scanners
type Registry struct {
	mu       sync.RWMutex
	scanners map[string]types.Scanner
}

// NewRegistry creates a new scanner registry
func NewRegistry() *Registry {
	return &Registry{
		scanners: make(map[string]types.Scanner),
	}
}

// Register adds a scanner to the registry
func (r *Registry) Register(scanner types.Scanner) error {
	if scanner == nil {
		return fmt.Errorf("scanner cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	scannerType := scanner.Type()
	if scannerType == "" {
		return fmt.Errorf("scanner type cannot be empty")
	}

	if _, exists := r.scanners[scannerType]; exists {
		log.Infof("Overriding existing scanner for type: %s", scannerType)
	}

	r.scanners[scannerType] = scanner
	log.Debugf("Registered scanner: %s (%s)", scanner.Name(), scannerType)

	return nil
}

// Get retrieves a scanner by type
func (r *Registry) Get(scannerType string) (types.Scanner, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	scanner, exists := r.scanners[scannerType]
	return scanner, exists
}

// List returns all registered scanner types
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]string, 0, len(r.scanners))
	for t := range r.scanners {
		types = append(types, t)
	}

	return types
}

// GetAllScanners returns all registered scanners
func (r *Registry) GetAllScanners() []types.Scanner {
	r.mu.RLock()
	defer r.mu.RUnlock()

	scanners := make([]types.Scanner, 0, len(r.scanners))
	for _, scanner := range r.scanners {
		scanners = append(scanners, scanner)
	}

	return scanners
}

// Count returns the number of registered scanners
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.scanners)
}

// Clear removes all registered scanners
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.scanners = make(map[string]types.Scanner)
	log.Debugf("Cleared all registered scanners")
}

// GetSuitableScanners returns scanners that can handle the given file
func (r *Registry) GetSuitableScanners(filePath string) []types.Scanner {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var suitable []types.Scanner

	for _, scanner := range r.scanners {
		if canHandle, err := scanner.CheckFile(filePath); err == nil && canHandle {
			suitable = append(suitable, scanner)
		}
	}

	return suitable
}

// Default registry instance
var defaultRegistry = NewRegistry()

// Register adds a scanner to the default registry
func Register(scanner types.Scanner) error {
	return defaultRegistry.Register(scanner)
}

// Get retrieves a scanner from the default registry
func Get(scannerType string) (types.Scanner, bool) {
	return defaultRegistry.Get(scannerType)
}

// List returns all scanner types from the default registry
func List() []string {
	return defaultRegistry.List()
}

// GetAllScanners returns all scanners from the default registry
func GetAllScanners() []types.Scanner {
	return defaultRegistry.GetAllScanners()
}

// GetSuitableScanners returns suitable scanners from the default registry
func GetSuitableScanners(filePath string) []types.Scanner {
	return defaultRegistry.GetSuitableScanners(filePath)
}

// Count returns the number of registered scanners in the default registry
func Count() int {
	return defaultRegistry.Count()
}
