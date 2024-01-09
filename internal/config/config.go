package config

import (
	"os"

	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/carbonetes/diggity/pkg/types"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

func Load() *types.Config {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	path := home + string(os.PathSeparator) + ".diggity.yaml"
	exist, err := helper.IsFileExists(path)
	if err != nil {
		return nil
	}

	if !exist {
		MakeDefaultConfigFile(path)
	}

	var config types.Config
	ReadConfigFile(&config, path)
	if config.Version != types.ConfigVersion {
		newConfig := New()
		err := mapstructure.Decode(config, &newConfig)
		if err != nil {
			log.Error(err)
		}
		newConfig.Version = types.ConfigVersion
		ReplaceConfigFile(newConfig, path)
	}
	stream.SetConfig(config)
	return &config
}

func New() types.Config {
	return types.Config{
		Version:      types.ConfigVersion,
		MaxFileSize:  10485760, // Set to 0 to disable limiter
		SecretConfig: LoadDefaultConfig(),
		AuthConfig: authn.AuthConfig{
			Username:      "",
			Password:      "",
			Auth:          "",
			IdentityToken: "",
			RegistryToken: "",
		},
		// Registry: types.RegistryParameters{
		// 	// URI:      "",
		// 	Username: "",
		// 	Password: "",
		// 	// Token:    "",
		// },
	}
}

func MakeDefaultConfigFile(path string) {
	os.Setenv("CONFIG_PATH", path)

	defaultConfig := New()
	err := helper.WriteYAML(defaultConfig, path)
	if err != nil {
		log.Error(err)
	}
}

func ReadConfigFile(config *types.Config, path string) {
	configFile, err := os.ReadFile(path)
	if err != nil {
		log.Error(err)
	}

	err = yaml.Unmarshal(configFile, config)
	if err != nil {
		log.Error(err)
	}
}

func ReplaceConfigFile(config types.Config, path string) {
	exist, err := helper.IsFileExists(path)
	if err != nil {
		log.Error(err)
	}

	if exist {
		err = os.Remove(path)
		if err != nil {
			log.Error(err)
		}
	}

	err = helper.WriteYAML(config, path)
	if err != nil {
		log.Error(err)
	}
}
