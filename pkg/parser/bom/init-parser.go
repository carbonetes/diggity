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

type ParserRequirements struct {
	Arguments *model.Arguments
	Dir       *string
	Contents  *[]model.Location
	Result    *model.Result
	WG        sync.WaitGroup
	Errors    *[]error
}

func NewParserRequirements(args *model.Arguments, dir *string, contents *[]model.Location) *ParserRequirements {
	return &ParserRequirements{
		Arguments: args,
		Dir:       dir,
		Contents:  contents,
		Errors:    new([]error),
		Result: &model.Result{
			Packages: new([]model.Package),
			Secret:   new(model.SecretResults),
			Distro:   new(model.Distro),
			SLSA:     new(model.SLSA),
		},
	}
}

// InitParsers initialize arguments
func InitParsers(argument model.Arguments) {
	Arguments = &argument
}
