package cmd

import (
	"os"

	"github.com/carbonetes/diggity/internal/config"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/internal/version"
	"github.com/carbonetes/diggity/pkg/reader"
	"github.com/carbonetes/diggity/pkg/scanner"
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
		Run: func(cmd *cobra.Command, args []string) {
			versionArg, _ := cmd.Flags().GetBool("version")
			if versionArg {
				log.Print(version.FromBuild().Version)
				os.Exit(0)
			}

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

			quiet, err := cmd.Flags().GetBool("quiet")
			if err != nil {
				log.Error(err.Error())
			}

			err = params.GetScanType()
			if err != nil {
				log.Error(err.Error())
			}

			outputFormat, err := cmd.Flags().GetString("output")
			if err != nil {
				log.Error(err.Error())
			}

			file, err := cmd.Flags().GetString("file")
			if err != nil {
				log.Error(err.Error())
			}

			scanners, err := cmd.Flags().GetStringArray("scanners")
			if err != nil {
				log.Error(err.Error())
			}

			if len(file) > 0 {
				params.SaveToFile = file
			}

			valid := types.IsValidOutputFormat(outputFormat)
			if !valid {
				log.Error("Invalid output format parameter")
			}

			params.Quiet = quiet
			params.SaveToFile = file
			params.Scanners = helper.SplitAndAppendStrings(scanners)
			params.OutputFormat = types.OutputFormat(outputFormat)
			params.AllowFileListing, err = cmd.Flags().GetBool("allow-file-listing")
			if err != nil {
				log.Error(err.Error())
			}
			reader.Init()
			stream.SetParameters(params)
		},
	}
)

func init() {
	config.Load()
	// Version sub command for indicating the version of diggity
	root.AddCommand(versionCmd)

	// Attest sub command for sbom attestation mechanism
	root.AddCommand(attestCmd)

	// Tarball flag to scan a tarball
	root.Flags().StringP("tar", "t", "", "Read a tarball from a path on disk for archives created from docker save (e.g. 'diggity path/to/image.tar)'")

	root.Flags().Bool("attest", false, "Add attestation to scan result")

	// Directory flag to scan a directory
	root.Flags().StringP("directory", "d", "", "Read directly from a path on disk (any directory) (e.g. 'diggity -fs path/to/directory)'")

	// Output flag to specify the output format
	root.Flags().StringP("output", "o", types.Table.String(), "Supported output types are: "+types.GetAllOutputFormat())

	// File flag to save the scan result to a file
	root.Flags().StringP("file", "f", "", "Save scan result to a file")

	// Quiet flag to allows the user to suppress all output except for errors
	root.Flags().BoolP("quiet", "q", false, "Suppress all output except for errors")

	// Scanners flag to specify the selected scanners to run
	root.Flags().StringArray("scanners", scanner.All, "Allow only selected scanners to run (e.g. --scanners apk,dpkg)")

	// File listing flag enables the user to list down all files related to the packages found
	root.Flags().Bool("allow-file-listing", false, "Allow parsers to list files related to the packages")

	// Version flag to print the version of diggity
	root.Flags().BoolP("version", "v", false, "Print the version of diggity")
}
