package config

import (
	"os"

	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/types"
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
		ReplaceConfigFile(New(), path)
	}

	return &config
}

func New() types.Config {
	return types.Config{
		Version:      types.ConfigVersion,
		MaxFileSize:  10485760,
		SecretConfig: LoadDefaultConfig(),
		Registry: types.RegistryParameters{
			URI:      "",
			Username: "",
			Password: "",
			Token:    "",
		},
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
