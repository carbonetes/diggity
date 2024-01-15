package cmd

import (
	"log"
	"os"

	"github.com/carbonetes/diggity/pkg/attest"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
	"github.com/spf13/cobra"
)

var (
	opts      types.AttestationOptions
	attestCmd = &cobra.Command{
		Use:   "attest",
		Short: "Attest generated SBOM.",
		Long:  "Generate and verify in-toto SBOM attesations with Cosign integrated with Diggity.",
		Args:  cobra.MaximumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
				os.Exit(0)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(opts.Predicate) == 0 {
				log.Fatal("No predicate has been given.")
			}
			config := stream.GetConfig()
			if len(opts.Key) == 0 {
				if len(config.Attestation.Key) == 0 {
					log.Fatal("No key detected.")
				}
				opts.Key = config.Attestation.Key
			}

			if len(opts.Pub) == 0 {
				if len(config.Attestation.Pub) == 0 {
					log.Fatal("No Pub detected.")
				}
				opts.Pub = config.Attestation.Pub
			}

			if len(opts.Password) == 0 {
				if len(config.Attestation.Password) == 0 {
					log.Fatal("No Password detected.")
				}
				opts.Password = config.Attestation.Password
			}
			attest.Run(args[0], opts)
		},
	}
)

func init() {
	attestCmd.Flags().StringVar(&opts.Key, "key", "", "Path to cosign.key used for the SBOM attestation")
	attestCmd.Flags().StringVar(&opts.Pub, "pub", "", "Path to cosign.pub used for the SBOM attestation")
	attestCmd.Flags().StringVar(&opts.AttestType, "type", "custom", "Type used for the attestation ([spdx, spdxjson, cdx, custom])")
	attestCmd.Flags().StringVar(&opts.Predicate, "predicate", "", "Path to the generated SBOM file to be attested")
	attestCmd.Flags().StringVar(&opts.Password, "password", "", "Password for the generated cosign key-pair")
	attestCmd.Flags().StringVarP(&opts.OutputFile, "output-file", "f", "", "Save the attestation result to the output file instead of writing to standard output")
	attestCmd.Flags().StringVarP(&opts.OutputType, "output", "o", "json", "Supported output types: \n[json, cyclonedx, cyclonedx-json, spdx-json, spdx-tag-value, github-json]")
}
