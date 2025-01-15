package config

import (
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func Load(configDirectory string, configFileName string, envPrefix string, configSt any) error {
	// if runtimeEnv == KUBERNETES {
	// 	mountPath := os.Getenv(MountConfigFilePath)
	// 	if mountPath == "" {
	// 		return errs.ErrArgs.WrapMsg(MountConfigFilePath + " env is empty")
	// 	}

	// 	return loadConfig(filepath.Join(mountPath, configFileName), envPrefix, config)
	// }

	return loadConfig(filepath.Join(configDirectory, configFileName), envPrefix, configSt)
}

func loadConfig(path string, envPrefix string, configSt any) error {
	viper.SetConfigFile(path)
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(configSt, func(config *mapstructure.DecoderConfig) {
		config.TagName = "mapstructure"
	}); err != nil {
		return err
	}
	return nil
}
