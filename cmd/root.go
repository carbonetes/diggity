package cmd

import (
	"os"

	"github.com/carbonetes/diggity/internal/curator"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/internal/presenter/status"
	"github.com/carbonetes/diggity/internal/scanner"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
	"github.com/spf13/cobra"
)

var (
	params = types.DefaultParameters()
	log    = logger.GetLogger()
	root   = &cobra.Command{
		Use:   "diggity",
		Args:  cobra.MaximumNArgs(1),
		Short: "BOM Diggity Scanner",
		Long:  `BOM Diggity is an open-source tool developed to streamline the critical process of generating a comprehensive Software Bill of Materials (SBOM) for Container Images and File Systems across various supported ecosystems.`,
		Run: func(cmd *cobra.Command, args []string) {
			tarball, _ := cmd.Flags().GetString("tar")
			directory, _ := cmd.Flags().GetString("directory")
			if len(args) > 0 {
				params.Input = helper.FormatImage(args[0])
			} else if len(tarball) > 0 {
				params.Input = tarball
			} else if len(directory) > 0 {
				params.Input = directory
			} else {
				_ = cmd.Help()
				os.Exit(0)
			}

			err := params.GetScanType()
			if err != nil {
				log.Error(err.Error())
			}

			outputFormat, err := cmd.Flags().GetString("output")
			if err != nil {
				log.Error(err.Error())
			}

			valid := types.IsValidOutputFormat(outputFormat)
			if !valid {
				log.Error("Invalid output format parameter")
			}

			params.OutputFormat = types.OutputFormat(outputFormat)
			params.AllowFileListing, err = cmd.Flags().GetBool("allow-file-listing")
			if err != nil {
				log.Error(err.Error())
			}
			params.AllowPullTimeout, err = cmd.Flags().GetBool("allow-pull-timeout")
			if err != nil {
				log.Error(err.Error())
			}
			stream.SetSecretParameters(params.Secrets)
			stream.SetParameters(params)

		},
	}
)

func init() {
	curator.Init()
	scanner.Init()
	status.Init()
}
