package config

type DsasConfig struct {
	IntegrationsDir string `yaml:"integrations_dir"`
	CoreApiConfig   struct {
		StartupHost int `yaml:"startup_port"`
	} `yaml:"core_api"`
}
