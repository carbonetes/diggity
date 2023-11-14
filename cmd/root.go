package cmd

import (
	"log"
	"os"

	"github.com/carbonetes/diggity/internal/cli"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/scanner"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
	"github.com/spf13/cobra"
)

var (
	params = types.DefaultParameters()
	root   = &cobra.Command{
		Use:   "diggity",
		Args:  cobra.MaximumNArgs(1),
		Short: "BOM Diggity Scanner",
		Long:  `BOM Diggity is an open-source tool developed to streamline the critical process of generating a comprehensive Software Bill of Materials (SBOM) for Container Images and File Systems across various supported ecosystems.`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
				os.Exit(0)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			tarball, _ := cmd.Flags().GetString("tar")
			directory, _ := cmd.Flags().GetString("directory")
			if len(args) > 0 {
				params.Input = helper.FormatImage(args[0])
			} else if len(tarball) > 0 {
				params.Input = tarball
			} else if len(directory) > 0 {
				params.Input = directory
			}

			err := params.GetScanType()
			if err != nil {
				log.Fatal(err.Error())
			}

			outputFormat, err := cmd.Flags().GetString("output")
			if err != nil {
				log.Fatal(err.Error())
			}

			valid := types.IsValidOutputFormat(outputFormat)
			if !valid {
				log.Fatal("Invalid output format parameter")
			}

			params.OutputFormat = types.OutputFormat(outputFormat)
			params.AllowFileListing, err = cmd.Flags().GetBool("allow-file-listing")
			if err != nil {
				log.Fatal(err.Error())
			}
			params.AllowPullTimeout, err = cmd.Flags().GetBool("allow-pull-timeout")
			if err != nil {
				log.Fatal(err.Error())
			}
			cli.Init()
			scanner.Init()
			stream.SetSecretParameters(params.Secrets)
			stream.SetParameters(params)

		},
	}
)
