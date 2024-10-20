package configs

type MainConfig struct {
	ConfigPath struct {
		DsasConfig string `yaml:"dsas_config_path"`
	} `yaml:"config_path"`
}
