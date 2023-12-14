package cmd

import (
	"github.com/carbonetes/diggity/internal/scanner"
	"github.com/carbonetes/diggity/pkg/types"
)

func init() {
	// Version sub command for indicating the version of diggity
	root.AddCommand(versionCmd)

	// Tarball flag to scan a tarball
	root.Flags().StringP("tar", "t", "", "Read a tarball from a path on disk for archives created from docker save (e.g. 'diggity path/to/image.tar)'")

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

	root.Flags().BoolP("version", "v", false, "Print the version of diggity")
}
