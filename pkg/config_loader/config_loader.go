package config_loader

import (
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

func LoadConfig[T any](
	configPath string,
	configStruct *T,
) (
	*T,
	error,
) {
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		slog.Error(
			"Failed to load config",
			"error",
			err,
		)
		return nil, err
	}
	err = yaml.Unmarshal(
		yamlFile,
		configStruct,
	)
	if err != nil {
		slog.Error(
			"Failed to load config",
			"error",
			err,
		)
	}

	return configStruct, err
}
