package cmd

import (
	"os"

	"github.com/carbonetes/diggity/internal/curator"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/logger"
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
			filesystem, _ := cmd.Flags().GetString("directory")
			if len(args) > 0 {
				params.Input = helper.FormatImage(args[0])
			} else if len(tarball) > 0 {
				params.Input = tarball
			} else if len(filesystem) > 0 {
				params.Input = filesystem
			} else {
				_ = cmd.Help()
				os.Exit(0)
			}

			err := params.GetScanType()
			if err != nil {
				log.Fatal(err.Error())
			}

			outputFormat, err := cmd.Flags().GetString("output")
			if err != nil {
				log.Fatal(err.Error())
			}

			file, err := cmd.Flags().GetString("file")
			if err != nil {
				log.Fatal(err.Error())
			}

			scanners, err := cmd.Flags().GetStringArray("scanners")
			if err != nil {
				log.Fatal(err.Error())
			}

			if len(file) > 0 {
				params.SaveToFile = file
			}

			valid := types.IsValidOutputFormat(outputFormat)
			if !valid {
				log.Fatal("Invalid output format parameter")
			}

			params.SaveToFile = file
			params.Scanners = helper.SplitAndAppendStrings(scanners)
			params.OutputFormat = types.OutputFormat(outputFormat)
			params.AllowFileListing, err = cmd.Flags().GetBool("allow-file-listing")
			if err != nil {
				log.Fatal(err.Error())
			}
			curator.Init()
			stream.SetParameters(params)
		},
	}
)
