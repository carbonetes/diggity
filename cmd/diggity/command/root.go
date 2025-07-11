package command

import (
	"context"

	analyzepkg "github.com/carbonetes/diggity/cmd/diggity/command/analyze"
	scanpkg "github.com/carbonetes/diggity/cmd/diggity/command/scan"
	versionpkg "github.com/carbonetes/diggity/cmd/diggity/command/version"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/config"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "diggity",
	Short: "BOM Diggity - Package, License, and Secret Detection Tool",
	Long: `BOM Diggity is an open-source tool for comprehensive software component analysis.
It performs scanning of raw data sources (images, filesystems) and analysis of structured 
SBOMs to detect packages, licenses, and secrets across various supported ecosystems.

Following industry-standard terminology:
- SCANNING: Examining raw container images, filesystems, and individual packages
- ANALYZING: Processing structured SBOMs (CycloneDX, SPDX) with rich contextual information`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.ExecuteContext(context.Background())
}

// Run executes the application
func Run() error {
	// Initialize configuration
	cfg := config.DefaultAnalysisConfig()
	scanCfg := config.DefaultScanConfig()

	log.Debugf("Diggity initialized with scan config: %+v", scanCfg)
	log.Debugf("Diggity initialized with analysis config: %+v", cfg)

	return Execute()
}

func init() {
	// Add global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose logging")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "suppress all output except errors")
	rootCmd.PersistentFlags().String("log-level", "info", "set log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().String("config", "", "config file path")

	// Add subcommands following jacked's pattern
	rootCmd.AddCommand(scanpkg.NewScanCommand())
	rootCmd.AddCommand(analyzepkg.NewAnalyzeCommand())
	rootCmd.AddCommand(versionpkg.NewVersionCommand())

	// Set default command to scan for backwards compatibility
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		// If no subcommand specified and we have arguments, run scan
		if len(args) > 0 {
			sc := scanpkg.NewScanCommand()
			sc.SetArgs(args)
			return sc.Execute()
		}
		return cmd.Help()
	}
}

// handleGlobalFlags processes global flags that affect logging and configuration
func handleGlobalFlags(cmd *cobra.Command) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	quiet, _ := cmd.Flags().GetBool("quiet")
	logLevel, _ := cmd.Flags().GetString("log-level")

	if verbose {
		log.SetLevel("debug")
	} else if quiet {
		log.SetLevel("error")
	} else if logLevel != "" {
		log.SetLevel(logLevel)
	}

	configPath, _ := cmd.Flags().GetString("config")
	if configPath != "" {
		log.Debugf("Using config file: %s", configPath)
		// Configuration loading would be implemented here
	}

	return nil
}
