package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/parser"
	versionPackage "github.com/carbonetes/diggity/internal/version"

	sbom "github.com/carbonetes/diggity/internal"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
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
	Arguments model.Arguments = model.Arguments{
		DisableFileListing:  new(bool),
		SecretContentRegex:  new(string),
		DisableSecretSearch: new(bool),
		Dir:                 new(string),
		Tar:                 new(string),
		Quiet:               new(bool),
		OutputFile:          new(string),
		ExcludedFilenames:   &[]string{},
		EnabledParsers:      &[]string{},
		RegistryURI:         new(string),
		RegistryUsername:    new(string),
		RegistryPassword:    new(string),
		RegistryToken:       new(string),
	}

	log = logger.GetLogger()

	diggity = &cobra.Command{
		Use:   "diggity",
		Args:  cobra.MaximumNArgs(1),
		Short: "BOM diggity SBOM Analyzer",
		Long:  `BOM Diggity's primary purpose is to ensure the security and integrity of software programs. It incorporates secret analysis allowing the user to secure crucial information before deploying any parts of the application to the public.`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 && !flagHasArg() {
				cmd.Help()
				os.Exit(0)
			}
			ValidateOutputArg(string(*Arguments.Output))
		},
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) > 0 {
				if flagHasArg() {
					log.Println(`"diggity [-d, --dir]" or diggity "[-t, --tar]" does not support with argument image`)
					os.Exit(127)
				}
				Arguments.Image = &args[0]
			} else if image != nil {
				Arguments.Image = image
			} else if flagHasArg() {
				//continue to sbom.Start
			} else {
				if len(args) == 0 || len(*Arguments.Image) == 0 {
					log.Printf(`"diggity [-i, --image]" is required or at least 1 argument "diggity [image]"`)
				}
				os.Exit(127)
			}
			sbom.Start(&Arguments)
		},
	}

	version = &cobra.Command{
		Use:   "version",
		Short: "Display Build Version Info Diggity",
		Long:  "Display Build Version Info Diggity",
		Args:  cobra.MaximumNArgs(0),
		RunE: func(_ *cobra.Command, _ []string) error {

			versionInfo := versionPackage.FromBuild()
			switch versionOutputFormat {
			case "text":
				// Version
				fmt.Println("Application:         ", versionInfo.AppName)
				fmt.Println("Version:             ", versionInfo.Version)
				fmt.Println("Build Date:          ", versionInfo.BuildDate)
				// Git
				fmt.Println("Git Commit:          ", versionInfo.GitCommit)
				fmt.Println("Git Description:     ", versionInfo.GitDesc)
				// Golang
				fmt.Println("Go Version:          ", versionInfo.GoVersion)
				fmt.Println("Compiler:            ", versionInfo.Compiler)
				fmt.Println("Platform:            ", versionInfo.Platform)
			case "json":

				jsonFormat := json.NewEncoder(os.Stdout)
				jsonFormat.SetEscapeHTML(false)
				jsonFormat.SetIndent("", " ")
				err := jsonFormat.Encode(&struct {
					model.Version
				}{
					Version: versionInfo,
				})
				if err != nil {
					return fmt.Errorf("show version information error: %+v", err)
				}
			default:
				return fmt.Errorf("unrecognize output format: %s", versionOutputFormat)
			}
			return nil
		},
	}

	config = &cobra.Command{
		Use:   "config",
		Short: "Display current configuration of diggity",
		Long:  "Display current configuration of diggity",
		Args:  cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			// Reset config file if specified
			if resetConfig {
				resetConfiguration()
				return
			}

			// Display contents of config file
			if displayConfig {
				config, _ := yaml.Marshal(DefaultConfig)
				log.Printf("%s", string(config))
				return
			}

			// Display path of config file
			if getConfigPath {
				log.Printf("Configuration directory: %s", os.Getenv("CONFIG_PATH"))
				return
			}

			// Show help
			cmd.Help()

		},
	}
)

func init() {
	intializeConfiguration()

	// diggity flags
	diggity.Flags().StringArrayVarP(outputArray, "output", "o", tableOutput, "Supported output types: \n[json, table, cyclonedx, cyclonedx-json, spdx-json, spdx-tag-value]")
	diggity.Flags().BoolVar(Arguments.DisableFileListing, "disable-file-listing", false, "Disables file listing from package metadata (default false)")
	diggity.Flags().StringVar(Arguments.SecretContentRegex, "secrets-content-regex", "", "Secret content regex are searched within files that matches the provided regular expression")
	diggity.Flags().BoolVar(Arguments.DisableSecretSearch, "disable-secret-search", false, "Disables secret search when set to true (default false)")
	diggity.Flags().BoolVarP(Arguments.Quiet, "quiet", "q", false, "Disable all output except SBOM result")
	diggity.Flags().StringVarP(Arguments.OutputFile, "output-file", "f", "", "Save the sbom result to the output file instead of writing to standard output")
	diggity.Flags().StringVarP(Arguments.Dir, "dir", "d", "", "Read directly from a path on disk (any directory) (e.g. 'diggity path/to/dir)'")
	diggity.Flags().StringVarP(Arguments.Tar, "tar", "t", "", "Read a tarball from a path on disk for archives created from docker save (e.g. 'diggity path/to/image.tar)'")
	diggity.Flags().Int64VarP(&Arguments.SecretMaxFileSize, "secret-max-file-size", "", 10485760, "Maximum file size that the secret will search -- each file")
	diggity.Flags().StringArrayVarP(Arguments.ExcludedFilenames, "secret-exclude-filenames", "", []string{}, "Exclude secret searching for each specified filenames")
	diggity.Flags().StringArrayVarP(Arguments.EnabledParsers, "enabled-parsers", "", []string{}, fmt.Sprintf("Specify enabled parsers (%+v) (default all)", parser.ParserNames))
	diggity.Flags().BoolVarP(&versionArg, "version", "v", false, "Display diggity version")
	diggity.Flags().StringVarP(Arguments.RegistryURI, "registry-uri", "", "index.docker.io/", "Registry uri endpoint")
	diggity.Flags().StringVarP(Arguments.RegistryUsername, "registry-username", "", "", "Username credential for private registry access")
	diggity.Flags().StringVarP(Arguments.RegistryPassword, "registry-password", "", "", "Password credential for private registry access")
	diggity.Flags().StringVarP(Arguments.RegistryToken, "registry-token", "", "", "Access token for private registry access")
	diggity.Flags().BoolVarP(&help, "help", "h", false, "Help for diggity")

	// version flags
	version.Flags().StringVarP(&versionOutputFormat, "output", "o", "text", "Format to display results ([text, json])")
	version.Flags().BoolVarP(&help, "help", "h", false, "Help for version")

	// config flags
	config.Flags().BoolVarP(&displayConfig, "display", "d", false, "Displays the contents of the configuration file")
	config.Flags().BoolVarP(&resetConfig, "reset", "r", false, "Restores default configuration file")
	config.Flags().BoolVarP(&getConfigPath, "path", "p", false, "Displays the path of the configuration file")
	config.Flags().BoolVarP(&help, "help", "h", false, "Help for config")

	cobra.OnInitialize(setPrioritizedArg)

	diggity.AddCommand(version)
	diggity.AddCommand(config)
	diggity.CompletionOptions.DisableDefaultCmd = true
}

// Initialize diggity yaml configuration
func intializeConfiguration() {

	home, err := os.UserHomeDir()
	if err != nil {
		err = errors.New("init-cmd: " + err.Error())
		parser.Errors = append(parser.Errors, &err)
	}

	ConfigDirectory = home + string(os.PathSeparator) + ".diggity.yaml"

	// Skip if file exists
	createConfiguration()

	os.Setenv("CONFIG_PATH", ConfigDirectory)
}

// Reset diggity config yaml file
func resetConfiguration() {
	err := os.Remove(ConfigDirectory)

	if err != nil {
		log.Print("[warning]: Unable to delete existing configuration file.")
	}

	createConfiguration()
	log.Println("Restored default configuration file.")
}

// Create diggity config yaml file
func createConfiguration() {
	if _, err := os.Stat(ConfigDirectory); errors.Is(err, os.ErrNotExist) {

		secretConfig := model.SecretConfig{
			Disabled:    false,
			SecretRegex: "API_KEY|SECRET_KEY|DOCKER_AUTH",
			Excludes:    &[]string{},
			MaxFileSize: 10485760,
		}
		DefaultConfig = model.Configuration{
			SecretConfig:   secretConfig,
			EnabledParsers: []string{},
			Output:         &[]string{},
			Quiet:          false,
			OutputFile:     "",
			Registry: model.Registry{
				URI:      "",
				Username: "",
				Password: "",
				Token:    "",
			},
		}

		yamlDefaultConfig, err := yaml.Marshal(&DefaultConfig)

		if err != nil {
			err = errors.New("init-cmd: " + err.Error())
			parser.Errors = append(parser.Errors, &err)
		}

		err = os.WriteFile(ConfigDirectory, yamlDefaultConfig, 0644)
		if err != nil {
			err = errors.New("init-cmd: " + err.Error())
			parser.Errors = append(parser.Errors, &err)
		}
	} else {
		// Read existing configuration instead
		configFile, _ := os.ReadFile(ConfigDirectory)
		err = yaml.Unmarshal(configFile, &DefaultConfig)
		if err != nil {
			panic(err)
		}
	}
}

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
}

// ValidateOutputArg checks if output types specified are valid
func ValidateOutputArg(outputType string) {
	for _, output := range strings.Split(outputType, ",") {
		if _, ok := model.OutputTypes[strings.ToLower(output)]; !ok {
			log.Printf("[warning]: Invalid output type: %+v \nSupported output types: [json, table, cyclonedx, cyclonedx-json, spdx-json, spdx-tag-value]", output)
			os.Exit(0)
		}
	}
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
