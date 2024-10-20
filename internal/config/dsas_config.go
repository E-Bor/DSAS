package config

type DsasConfig struct {
	IntegrationsDir string `yaml:"integrations_dir"`
	CoreApiConfig   struct {
		StartupHost string `yaml:"startup_host"`
	} `yaml:"core_api"`
}
