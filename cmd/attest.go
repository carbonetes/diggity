package cmd

import (
	"os"

	"github.com/carbonetes/diggity/pkg/attestation"
	"github.com/spf13/cobra"
)

var (
	attest = &cobra.Command{
		Use:   "attest",
		Short: "Attest generated SBOM.",
		Long:  "Generate and verify in-toto SBOM attesations with Cosign integrated with Diggity.",
		Args:  cobra.MaximumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 && !flagHasArg() {
				_ = cmd.Help()
				os.Exit(0)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			attestation.Attest(args[0], &attestationOptions)
		},
	}
)
