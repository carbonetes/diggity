package scan

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/config"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/scan"
	"github.com/spf13/cobra"
)

// NewScanCommand creates the scan command following the terminology refactoring
func NewScanCommand() *cobra.Command {
	var (
		targetPath      string
		outputFormat    string
		outputFile      string
		includeSecrets  bool
		includeLicenses bool
		quiet           bool
		configFile      string
	)

	cmd := &cobra.Command{
		Use:   "scan [TARGET]",
		Short: "Scan raw data sources for packages, licenses, and secrets",
		Long: `Scan performs discovery and detection in unstructured data sources including:
- Container images
- Filesystems and directories  
- Individual packages
- Archive files

This follows the terminology refactoring: "Scanning" for examining raw data sources.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runScan(cmd.Context(), args, scanOptions{
				targetPath:      targetPath,
				outputFormat:    outputFormat,
				outputFile:      outputFile,
				includeSecrets:  includeSecrets,
				includeLicenses: includeLicenses,
				quiet:           quiet,
				configFile:      configFile,
			})
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&targetPath, "target", "t", "", "target to scan (image, directory, file)")
	cmd.Flags().StringVarP(&outputFormat, "format", "f", "json", "output format (json, cyclonedx, spdx)")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "output file path (default: stdout)")
	cmd.Flags().BoolVar(&includeSecrets, "secrets", true, "include secret scanning")
	cmd.Flags().BoolVar(&includeLicenses, "licenses", true, "include license scanning")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "suppress output except errors")
	cmd.Flags().StringVar(&configFile, "config", "", "config file path")

	// Legacy flags for backward compatibility
	cmd.Flags().String("dir", "", "scan directory (deprecated: use --target)")
	cmd.Flags().String("tar", "", "scan tarball (deprecated: use --target)")
	cmd.Flags().String("image", "", "scan container image (deprecated: use --target)")

	return cmd
}

type scanOptions struct {
	targetPath      string
	outputFormat    string
	outputFile      string
	includeSecrets  bool
	includeLicenses bool
	quiet           bool
	configFile      string
}

func runScan(ctx context.Context, args []string, opts scanOptions) error {
	// Determine target from args or flags
	target := determineTarget(args, opts)
	if target == "" {
		return fmt.Errorf("no target specified. Use --target or provide as argument")
	}

	if opts.quiet {
		log.SetLevel("error")
	}

	log.Infof("Starting scan of target: %s", target)
	start := time.Now()

	// Load configuration
	cfg := config.DefaultScanConfig()
	if opts.configFile != "" {
		log.Debugf("Loading config from: %s", opts.configFile)
		// Configuration loading would be implemented here
	}

	// Configure scanning options
	cfg.SecretScanningEnabled = opts.includeSecrets
	cfg.LicenseScanningEnabled = opts.includeLicenses

	// Create scan engine with all 16 parsers
	engine := scan.NewEngine(cfg)

	// Show engine status
	status := engine.GetEngineStatus()
	log.Infof("Scan engine ready with %d parsers (%d language, %d system package, %d other)",
		status["total_parsers"],
		status["programming_languages"],
		status["system_packages"],
		status["other_parsers"])

	if !opts.quiet {
		log.Debugf("Available parsers: %v", engine.GetAvailableScanners())
	}

	// Create scan target
	scanTarget := scan.ScanTarget{
		Type: determineScanTargetType(target),
		Path: target,
	}

	// Perform scan
	result, err := engine.ScanTarget(ctx, scanTarget)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	duration := time.Since(start)
	log.Infof("Scan completed in %v", duration)

	// Output results
	return outputResults(result, opts.outputFormat, opts.outputFile)
}

func determineTarget(args []string, opts scanOptions) string {
	// Check args first
	if len(args) > 0 {
		return args[0]
	}

	// Check explicit target flag
	if opts.targetPath != "" {
		return opts.targetPath
	}

	// Check legacy flags for backward compatibility
	// These would be retrieved from the cobra command flags
	return ""
}

func determineScanTargetType(target string) scan.TargetType {
	// Determine target type based on path characteristics
	if strings.Contains(target, ":") && !filepath.IsAbs(target) {
		// Likely a container image reference
		return scan.TargetTypeImage
	}

	if strings.HasSuffix(target, ".tar") || strings.HasSuffix(target, ".tar.gz") {
		return scan.TargetTypeArchive
	}

	// Check if it's a file or directory
	if info, err := os.Stat(target); err == nil {
		if info.IsDir() {
			return scan.TargetTypeDirectory
		}
		return scan.TargetTypeFile
	}

	// Default to directory for paths that don't exist yet
	return scan.TargetTypeDirectory
}

func outputResults(result interface{}, format, outputFile string) error {
	// Output formatting and writing would be implemented here
	log.Infof("Results would be output in %s format", format)
	if outputFile != "" {
		log.Infof("Output would be written to: %s", outputFile)
	}

	// Show some information about the integrated scanners
	if compResult, ok := result.(*model.ComprehensiveScanResult); ok {
		log.Infof("Scan completed with %d errors", len(compResult.Errors))
		if compResult.PackageResult != nil {
			log.Infof("Found %d packages", compResult.PackageResult.ComponentsFound)
		}
		if compResult.SecretResult != nil {
			log.Infof("Found %d secrets", compResult.SecretResult.SecretsFound)
		}
		if compResult.LicenseResult != nil {
			log.Infof("Found %d licenses", compResult.LicenseResult.LicensesFound)
		}
	}

	log.Info("Scan results processed successfully")
	return nil
}
