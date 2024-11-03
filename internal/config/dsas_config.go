package config

type DsasConfig struct {
	IntegrationsDir string `yaml:"integrations_dir"`
	CoreApiConfig   struct {
		StartupHost int `yaml:"startup_port"`
	} `yaml:"core_api"`
	DsasCoreConfig struct {
		LoadWorkersCount            int   `yaml:"load_workers_count"`
		DefaultAverageLoadTime      int64 `yaml:"default_average_load_time"` // seconds
		TraceIdLength               int   `yaml:"trace_id_length"`
		QueueSleepTime              int   `yaml:"queue_sleep_time"` // seconds
		QueueLength                 int   `yaml:"queue_length"`
		WorkerPoolChannelBufferSize int   `yaml:"worker_pool_channel_buffer_size"`
		WorkerSleepTime             int   `yaml:"worker_sleep_time"`
	} `yaml:"dsas_core_config"`
}
