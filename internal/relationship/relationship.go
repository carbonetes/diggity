package relationship

import (
	"sync"

	"github.com/carbonetes/diggity/pkg/model"
)

type (
	relationshipFinder []func()
)

var (
	// Relationships - common collection of relationships
	Relationships []model.Relationship

	// FindRelationships - collection of the relationshipFinderFunctions
	FindRelationships = relationshipFinder{
		FindSourceContains,
	}

	wg sync.WaitGroup
)

func GetRelationships(results map[string]*model.Package) []model.Relationship {
	wg.Add(len(FindRelationships))
	for _, rf := range FindRelationships {
		go rf()
	}
	wg.Wait()
	return Relationships
}
