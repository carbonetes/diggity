package cli

import (
	"github.com/carbonetes/diggity/internal/curator"
	"github.com/carbonetes/diggity/internal/presenter/status"
	"github.com/carbonetes/diggity/internal/scanner"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
)

func init() {
	curator.Init()
	scanner.Init()
	status.Init()
}

func Start(params types.Parameters) {
	stream.SetParameters(params)
}
