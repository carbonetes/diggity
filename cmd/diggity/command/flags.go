package command

import "github.com/carbonetes/diggity/pkg/types"

func init() {
	// Version sub command for indicating the version of diggity
	scan.AddCommand(versionCmd)

	// Attest sub command for sbom attestation mechanism
	scan.AddCommand(attestCmd)

	// Tarball flag to scan a tarball
	scan.Flags().StringP("tar", "t", "", "Read a tarball from a path on disk for archives created from docker save (e.g. 'diggity -t path/to/image.tar)'")

	scan.Flags().Bool("attest", false, "Add attestation to scan result")

	// Directory flag to scan a directory
	scan.Flags().StringP("dir", "d", "", "Read directly from a path on disk (any directory) (e.g. 'diggity -d path/to/directory)'")

	// Output flag to specify the output format
	scan.Flags().StringP("output", "o", types.Table.String(), "Supported output types are: "+types.GetAllOutputFormat())

	// File flag to save the scan result to a file
	scan.Flags().StringP("file", "f", "", "Save scan result to a file")

	// Quiet flag to allows the user to suppress all output except for errors
	scan.Flags().BoolP("quiet", "q", false, "Suppress all output except for errors")

	// Version flag to print the version of diggity
	scan.Flags().BoolP("version", "v", false, "Print the version of diggity")
}