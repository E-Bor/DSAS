package configs

import "DSAS/pkg/dsas_logger"

type MainConfig struct {
	Env        dsas_logger.LogLvl `yaml:"env"`
	ConfigPath struct {
		DsasConfig string `yaml:"dsas_config_path"`
	} `yaml:"config_path"`
}
