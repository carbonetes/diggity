package cmd

import (
	"errors"
	"os"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
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
			_ = cmd.Help()

		},
	}
)

// Initialize diggity yaml configuration
func intializeConfiguration() {

	home, err := os.UserHomeDir()
	if err != nil {
		err = errors.New("init-cmd: " + err.Error())
		log.Fatal(err)
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
			Extensions:  &model.DefaultSecretExtensions,
		}
		attestationConfig := model.AttestationConfig{
			Key:      "cosign.key",
			Pub:      "cosign.pub",
			Password: "",
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
			AttestationConfig: attestationConfig,
		}

		yamlDefaultConfig, err := yaml.Marshal(&DefaultConfig)

		if err != nil {
			err = errors.New("init-cmd: " + err.Error())
			log.Fatal(err)
		}

		err = os.WriteFile(ConfigDirectory, yamlDefaultConfig, 0644)
		if err != nil {
			err = errors.New("init-cmd: " + err.Error())
			log.Fatal(err)
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
