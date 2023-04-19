package bom

import (
	"sync"

	"github.com/carbonetes/diggity/pkg/model"
)

var (
	Target *string
	// Arguments - CLI Arguments
	Arguments *model.Arguments
	// Packages - common collection of packages found by parsers
	Packages []*model.Package
	// WG - common waitgroup for all the parsers
	WG sync.WaitGroup
	// Errors - common errors encountered by parsers
	Errors []*error
)

// InitParsers initialize arguments
func InitParsers(argument model.Arguments) {
	Arguments = &argument
}
