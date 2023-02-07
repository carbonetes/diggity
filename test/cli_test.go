package test

import (
	"testing"

	"github.com/carbonetes/diggity/internal/model"

	"github.com/spf13/cobra"
)

var (
	t         *testing.T
	Arguments model.Arguments

	// Test Variables
	argsImage = "alpine"
)

var rootCmd = &cobra.Command{
	Use:    "diggity [image] [flags]",
	Short:  "BOM diggity SBOM Analyzer",
	Long:   `BOM Diggity's primary purpose is to ensure the security and integrity of software programs. It incorporates secret analysis allowing the user to secure crucial information before deploying any parts of the application to the public.`,
	PreRun: preRun,
}

func TestCli(t *testing.T) {
	Execute()
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	err := rootCmd.Execute()
	if err != nil {
		t.Fail()
	}
}

func preRun(_ *cobra.Command, args []string) {
	args = append(args, argsImage)
	if len(args) == 0 {
		t.Fail()
	}
}
