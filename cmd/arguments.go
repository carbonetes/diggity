package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/carbonetes/diggity/internal/logger"
	versionPackage "github.com/carbonetes/diggity/internal/version"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/util"
	"github.com/spf13/cobra"
)

var (
	// Output Type Variables
	defaultOutput = "table"
	tableOutput   = []string{"table"}
	outputArray   = new([]string)
	// DefaultConfig = Marshalled from file in ConfigDirectory
	DefaultConfig model.Configuration
	// ConfigDirectory = $HOME/.diggity.yaml
	ConfigDirectory     string
	displayConfig       bool
	resetConfig         bool
	getConfigPath       bool
	help                bool // help flag
	image               *string
	versionArg          bool
	versionOutputFormat string
	// Arguments is an instance of the actual arguments passed
	Arguments = model.NewArguments()

	log = logger.GetLogger()

	// attestation
	attestationOptions = model.AttestationOptions{
		Key:        new(string),
		Pub:        new(string),
		AttestType: new(string),
		Predicate:  new(string),
		Password:   new(string),
		OutputFile: new(string),
		OutputType: new(string),
		BomArgs:    new(model.Arguments),
		Provenance: new(string),
	}
)

func init() {
	intializeConfiguration()

	// diggity flags
	diggity.Flags().StringArrayVarP(outputArray, "output", "o", tableOutput, fmt.Sprintf("Supported output types: \n%+v", model.OutputList))
	diggity.Flags().BoolVar(Arguments.DisableFileListing, "disable-file-listing", false, "Disables file listing from package metadata (default false)")
	diggity.Flags().BoolVar(Arguments.DisablePullTimeout, "disable-pull-timeout", false, "Disables the timeout when pulling an image from server (default false)")
	diggity.Flags().StringVar(Arguments.SecretContentRegex, "secrets-content-regex", "", "Secret content regex are searched within files that matches the provided regular expression")
	diggity.Flags().BoolVar(Arguments.DisableSecretSearch, "disable-secret-search", false, "Disables secret search when set to true (default false)")
	diggity.Flags().BoolVarP(Arguments.Quiet, "quiet", "q", false, "Disable all output except SBOM result")
	diggity.Flags().StringVarP(Arguments.OutputFile, "output-file", "f", "", "Save the sbom result to the output file instead of writing to standard output")
	diggity.Flags().StringVarP(Arguments.Dir, "dir", "d", "", "Read directly from a path on disk (any directory) (e.g. 'diggity path/to/dir)'")
	diggity.Flags().StringVarP(Arguments.Tar, "tar", "t", "", "Read a tarball from a path on disk for archives created from docker save (e.g. 'diggity path/to/image.tar)'")
	diggity.Flags().Int64VarP(&Arguments.SecretMaxFileSize, "secret-max-file-size", "", 10485760, "Maximum file size that the secret will search -- each file")
	diggity.Flags().StringArrayVarP(Arguments.ExcludedFilenames, "secret-exclude-filenames", "", []string{}, "Exclude secret searching for each specified filenames")
	diggity.Flags().StringArrayVarP(Arguments.SecretExtensions, "secret-extensions", "", []string{}, "File extensions to consider for secret search (default no extension)")
	diggity.Flags().StringArrayVarP(Arguments.EnabledParsers, "enabled-parsers", "", []string{}, fmt.Sprintf("Specify enabled parsers (%+v) (default all)", util.ParserNames))
	diggity.Flags().BoolVarP(&versionArg, "version", "v", false, "Display diggity version")
	diggity.Flags().StringVarP(Arguments.RegistryURI, "registry-uri", "", "index.docker.io/", "Registry uri endpoint")
	diggity.Flags().StringVarP(Arguments.RegistryUsername, "registry-username", "", "", "Username credential for private registry access")
	diggity.Flags().StringVarP(Arguments.RegistryPassword, "registry-password", "", "", "Password credential for private registry access")
	diggity.Flags().StringVarP(Arguments.RegistryToken, "registry-token", "", "", "Access token for private registry access")
	diggity.Flags().StringVarP(Arguments.Provenance, "provenance", "", "", "Provenance file to include in the SBOM")
	diggity.Flags().BoolVarP(&help, "help", "h", false, "Help for diggity")

	// version flags
	version.Flags().StringVarP(&versionOutputFormat, "output", "o", "text", "Format to display results ([text, json])")
	version.Flags().BoolVarP(&help, "help", "h", false, "Help for version")

	// config flags
	config.Flags().BoolVarP(&displayConfig, "display", "d", false, "Displays the contents of the configuration file")
	config.Flags().BoolVarP(&resetConfig, "reset", "r", false, "Restores default configuration file")
	config.Flags().BoolVarP(&getConfigPath, "path", "p", false, "Displays the path of the configuration file")
	config.Flags().BoolVarP(&help, "help", "h", false, "Help for configuration")

	// attest flags
	attest.Flags().StringVarP(attestationOptions.Key, "key", "k", "", "Path to cosign.key used for the SBOM attestation")
	attest.Flags().StringVarP(attestationOptions.Pub, "pub", "p", "", "Path to cosign.pub used for the SBOM attestation")
	attest.Flags().StringVarP(attestationOptions.AttestType, "type", "t", "custom", "Type used for the attestation ([spdx, spdxjson, cyclonedx, custom])")
	attest.Flags().StringVar(attestationOptions.Predicate, "predicate", "", "Path to the generated SBOM file to be attested")
	attest.Flags().StringVar(attestationOptions.Password, "password", "", "Password for the generated cosign key-pair")
	attest.Flags().StringVarP(attestationOptions.OutputFile, "output-file", "f", "", "Save the attestation result to the output file instead of writing to standard output")
	attest.Flags().StringVarP(attestationOptions.OutputType, "output", "o", "json", "Supported output types: \n[json, cyclonedx, cyclonedx-json, spdx-json, spdx-tag-value, github-json]")
	attest.Flags().StringVar(attestationOptions.Provenance, "provenance", "", "Provenance file to include in the SBOM")
	attest.Flags().BoolVarP(&help, "help", "h", false, "Help for attest")

	cobra.OnInitialize(setPrioritizedArg)
	cobra.OnInitialize(setAttestArgs)
	attestationOptions.BomArgs = Arguments

	diggity.AddCommand(version)
	diggity.AddCommand(config)
	diggity.AddCommand(attest)
	diggity.CompletionOptions.DisableDefaultCmd = true
}

// TODO: reduce code complexity

// Define either model.Argument or model.Configuration will be prioritized
func setPrioritizedArg() {
	if versionArg {
		log.Printf("diggity %s", versionPackage.FromBuild().Version)
		os.Exit(0)
	}

	// Set values from flags, if applicable
	setArrayArgs()

	// Check if flags where specified, else use config values
	if !diggity.Flags().Lookup("output").Changed {
		if len(*DefaultConfig.Output) > 0 {
			defaultConfigOutput := strings.Join(*DefaultConfig.Output, ",")
			Arguments.Output = (*model.Output)(&defaultConfigOutput)
		} else {
			// Display table output if there is no specification in the config file
			Arguments.Output = (*model.Output)(&defaultOutput)
		}
	}
	if !diggity.Flags().Lookup("enabled-parsers").Changed {
		Arguments.EnabledParsers = &DefaultConfig.EnabledParsers
	}
	if !diggity.Flags().Lookup("disable-file-listing").Changed && !*Arguments.DisableFileListing {
		Arguments.DisableFileListing = &DefaultConfig.DisableFileListing
	}
	if !diggity.Flags().Lookup("disable-pull-timeout").Changed && !*Arguments.DisablePullTimeout {
		Arguments.DisablePullTimeout = &DefaultConfig.DisablePullTimeout
	}
	if !diggity.Flags().Lookup("secrets-content-regex").Changed {
		Arguments.SecretContentRegex = &DefaultConfig.SecretConfig.SecretRegex
	}
	if !diggity.Flags().Lookup("disable-secret-search").Changed && !*Arguments.DisableSecretSearch {
		Arguments.DisableSecretSearch = &DefaultConfig.SecretConfig.Disabled
	}
	if !diggity.Flags().Lookup("quiet").Changed && !*Arguments.Quiet {
		Arguments.Quiet = &DefaultConfig.Quiet
	}
	if !diggity.Flags().Lookup("output-file").Changed {
		Arguments.OutputFile = &DefaultConfig.OutputFile
	}
	if !diggity.Flags().Lookup("secret-max-file-size").Changed {
		Arguments.SecretMaxFileSize = DefaultConfig.SecretConfig.MaxFileSize
	}
	if !diggity.Flags().Lookup("secret-exclude-filenames").Changed {
		Arguments.ExcludedFilenames = DefaultConfig.SecretConfig.Excludes
	}
	if !diggity.Flags().Lookup("secret-extensions").Changed {
		Arguments.SecretExtensions = DefaultConfig.SecretConfig.Extensions
	}
	if !diggity.Flags().Lookup("registry-uri").Changed {
		Arguments.RegistryURI = &DefaultConfig.Registry.URI
	}
	if !diggity.Flags().Lookup("registry-username").Changed {
		Arguments.RegistryUsername = &DefaultConfig.Registry.Username
	}
	if !diggity.Flags().Lookup("registry-password").Changed {
		Arguments.RegistryPassword = &DefaultConfig.Registry.Password
	}
	if !diggity.Flags().Lookup("registry-token").Changed {
		Arguments.RegistryToken = &DefaultConfig.Registry.Token
	}
}

// Define args for attest
func setAttestArgs() {
	attestationConfig := DefaultConfig.AttestationConfig
	if !attest.Flags().Lookup("key").Changed {
		*attestationOptions.Key = attestationConfig.Key
	}
	if !attest.Flags().Lookup("pub").Changed && !*Arguments.DisableFileListing {
		*attestationOptions.Pub = attestationConfig.Pub
	}
	if !attest.Flags().Lookup("password").Changed {
		*attestationOptions.Password = attestationConfig.Password
	}
	if attest.Flags().Lookup("output").Changed {
		ValidateOutputArg(*attestationOptions.OutputType)
	}
}

// Check if flag has dir or tar arguments
func flagHasArg() bool {
	if len(*Arguments.Dir) != 0 || len(*Arguments.Tar) != 0 {
		return true
	}
	return false
}

// Set values form flags with string array var
func setArrayArgs() {
	// Set output from flags
	output := strings.Join(*outputArray, ",")
	Arguments.Output = (*model.Output)(&output)

	// Set enabled parsers from flags
	enabledParsers := SplitArgs(*Arguments.EnabledParsers)
	Arguments.EnabledParsers = &enabledParsers

	// Set excluded filenames from flags
	excludedFilenames := SplitArgs(*Arguments.ExcludedFilenames)
	Arguments.ExcludedFilenames = &excludedFilenames

	// Set secret extensions from flags
	secretExtensions := SplitArgs(*Arguments.SecretExtensions)
	Arguments.SecretExtensions = &secretExtensions
}

// ValidateOutputArg checks if output types specified are valid
func ValidateOutputArg(outputType string) error {
	for _, output := range strings.Split(outputType, ",") {
		// Validate from Default Output Types
		if _, ok := model.OutputTypes[strings.ToLower(output)]; !ok {
			return fmt.Errorf("Invalid output type: %+v \nSupported output types: %+v", output, model.OutputList)
		}
	}
	return nil
}

// SplitArgs splits arguments with comma, if any
func SplitArgs(args []string) (result []string) {
	for _, arg := range args {
		if !strings.Contains(arg, ",") {
			result = append(result, arg)
			continue
		}
		result = append(result, strings.Split(arg, ",")...)
	}
	return result
}
