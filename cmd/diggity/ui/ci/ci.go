package ci

import (
	"os"
	"time"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/ci/api"
	stream "github.com/carbonetes/diggity/cmd/diggity/grove"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/scanner/secret"
	"github.com/carbonetes/diggity/pkg/types"
	"github.com/golistic/urn"
)

func Run(result *cyclonedx.BOM, addr *urn.URN, params types.Parameters, duration float64) {
	// Carbonetes CI API
	start := time.Now().Add(-time.Duration(duration * float64(time.Second)))
	api.SavePluginRepository(result, params.Input, params.Plugin, start, 1, 0)
	// Secrets
	evaluateSecrets(addr)
	// Components
	showComponentsResult(result)
	os.Exit(0)
}

func evaluateSecrets(addr *urn.URN) {
	cdx.New(addr)
	secretAddr := *addr
	secretAddr.NID = "secret"
	secret.New(&secretAddr)
	s, _ := stream.Get(secretAddr.String())
	secrets := s.([]types.Secret)
	if len(secrets) == 0 {
		log.Printf("Assessment Passed: No Secrets Found.")
	} else {
		log.Printf("Total Secrets Found: %v", len(secrets))
		SecretTable(secrets)
		log.Printf("Assessment Failed: Secrets Found!")
	}
}

func showComponentsResult(result *cyclonedx.BOM) {
	totalComponents := len(*result.Components)
	if totalComponents == 0 {
		log.Printf("No Packages Found.")
	} else {
		log.Printf("\nTotal Packages Found: %v", totalComponents)
		Table(result)
	}
}
