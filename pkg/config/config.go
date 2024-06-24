package config

import (
	"os"

	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/types"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

var Config types.Config

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Debug(err)
	}
	path := home + string(os.PathSeparator) + ".diggity.yaml"
	exist, err := helper.IsFileExists(path)
	if err != nil {
		log.Debug(err)
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
			log.Debug(err)
		}
		newConfig.Version = types.ConfigVersion
		ReplaceConfigFile(newConfig, path)
	}
	Config = config
}

func New() types.Config {
	return types.Config{
		Version:      types.ConfigVersion,
		MaxFileSize:  52428800, // 50MB
		SecretConfig: LoadDefaultConfig(),
		Registry: types.RegistryConfig{
			Username: "",
			Password: "",
		},
	}
}

func MakeDefaultConfigFile(path string) {
	os.Setenv("CONFIG_PATH", path)

	defaultConfig := New()
	err := helper.WriteYAML(defaultConfig, path)
	if err != nil {
		log.Debug(err)
	}
}

func ReadConfigFile(config *types.Config, path string) {
	configFile, err := os.ReadFile(path)
	if err != nil {
		log.Debug(err)
	}

	err = yaml.Unmarshal(configFile, config)
	if err != nil {
		log.Debug(err)
	}
}

func ReplaceConfigFile(config types.Config, path string) {
	exist, err := helper.IsFileExists(path)
	if err != nil {
		log.Debug(err)
	}

	if exist {
		err = os.Remove(path)
		if err != nil {
			log.Debug(err)
		}
	}

	err = helper.WriteYAML(config, path)
	if err != nil {
		log.Debug(err)
	}
}
