package cmd

import (
	"os"

	"github.com/carbonetes/diggity/internal/cli"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/internal/version"
	"github.com/carbonetes/diggity/pkg/types"
	"github.com/spf13/cobra"
)

var (
	root = &cobra.Command{
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
			filesystem, _ := cmd.Flags().GetString("dir")

			params := types.DefaultParameters()
			if len(args) > 0 {
				params.Input = helper.FormatImage(args[0])
			} else if len(tarball) > 0 {
				params.Input = tarball
			} else if len(filesystem) > 0 {
				if found, _ := helper.IsDirExists(filesystem); !found {
					log.Fatal("directory not found: " + filesystem)
					return
				}
				params.Input = filesystem
			} else {
				_ = cmd.Help()
				os.Exit(0)
			}

			quiet, err := cmd.Flags().GetBool("quiet")
			if err != nil {
				log.Debug(err.Error())
			}

			err = params.GetScanType()
			if err != nil {
				log.Debug(err.Error())
			}

			outputFormat, err := cmd.Flags().GetString("output")
			if err != nil {
				log.Debug(err.Error())
			}

			file, err := cmd.Flags().GetString("file")
			if err != nil {
				log.Debug(err.Error())
			}

			scanners, err := cmd.Flags().GetStringArray("scanners")
			if err != nil {
				log.Debug(err.Error())
			}

			if len(file) > 0 {
				params.SaveToFile = file
			}

			valid := types.IsValidOutputFormat(outputFormat)
			if !valid {
				log.Debug("Invalid output format parameter")
			}

			params.Quiet = quiet
			params.SaveToFile = file
			params.Scanners = helper.SplitAndAppendStrings(scanners)
			params.OutputFormat = types.OutputFormat(outputFormat)
			if err != nil {
				log.Debug(err.Error())
			}

			// CI Mode
			params.CI, _ = cmd.Flags().GetBool("ci")
			params.Token, _ = cmd.Flags().GetString("token")
			params.Plugin, _ = cmd.Flags().GetString("plugin")

			// If CI mode is enabled, suppress all output except for errors
			if params.CI {
				params.Quiet = true
				if params.Token == "" {
					log.Fatal("Token is required. Please generate a token at https://app.carbonetes.com/personal-access-token and set it using the --token flag.")
					os.Exit(1)
				}
			}
			if params.Quiet {
				params.OutputFormat = types.OutputFormat("json")
			}
			if !(len(params.Plugin) > 0) {
				params.Plugin = "oss"
			}

			cli.Start(params)
		},
	}
)

func init() {
	// Version sub command for indicating the version of diggity
	root.AddCommand(versionCmd)

	// Attest sub command for sbom attestation mechanism
	root.AddCommand(attestCmd)

	// Tarball flag to scan a tarball
	root.Flags().StringP("tar", "t", "", "Read a tarball from a path on disk for archives created from docker save (e.g. 'diggity -t path/to/image.tar)'")

	root.Flags().Bool("attest", false, "Add attestation to scan result")

	// Directory flag to scan a directory
	root.Flags().StringP("dir", "d", "", "Read directly from a path on disk (any directory) (e.g. 'diggity -d path/to/directory)'")

	// Output flag to specify the output format
	root.Flags().StringP("output", "o", types.Table.String(), "Supported output types are: "+types.GetAllOutputFormat())

	// File flag to save the scan result to a file
	root.Flags().StringP("file", "f", "", "Save scan result to a file")

	// Quiet flag to allows the user to suppress all output except for errors
	root.Flags().BoolP("quiet", "q", false, "Suppress all output except for errors")

	// Scanners flag to specify the selected scanners to run
	// root.Flags().StringArray("scanners", scanner.All, "Allow only selected scanners to run (e.g. --scanners apk,dpkg)")

	// Version flag to print the version of diggity
	root.Flags().BoolP("version", "v", false, "Print the version of diggity")

	// CI flag to enable CI mode
	// CI mode is a mode that is used to run jacked in a CI/CD pipeline
	root.Flags().BoolP("ci", "", false, "Enable CI mode [experimental] (e.g. --ci)")
	root.Flags().StringP("token", "", "", "CI mode requires a personal access token. Sign up at https://app.carbonetes.com/ and generate your token to enable integration.")
	root.Flags().StringP("plugin", "", "", "CI mode set plugin type. (default jacked)")
}
