package config

type AppConfig struct {
	Server struct {
		Port          string `mapstructure:"port"`
		CertFile      string `mapstructure:"cert_file"`
		KeyFile       string `mapstructure:"key_file"`
		DockerHostURL string `mapstructure:"docker_host_url"`
	} `mapstructure:"server"`
	Database struct {
		ConnectionURI string `mapstructure:"connection_uri"`
		DbName        string `mapstructure:"db_name"`
	} `mapstructure:"database"`
	GameServer struct {
		HostingDir string   `mapstructure:"hosting_dir"`
		WorldDir   string   `mapstructure:"world_dir"`
		ImageName  string   `mapstructure:"image_name"`
		EnvVars    []string `mapstructure:"env_vars"`
		Volumes    []string `mapstructure:"volumes"`
		Labels     []string `mapstructure:"labels"`
		Networks   []string `mapstructure:"networks"`
		Logging    struct {
			Driver  string            `mapstructure:"driver"`
			Options map[string]string `mapstructure:"options"`
		}
	} `mapstructure:"game_server"`
}

func (a AppConfig) UseTLS() bool {
	return a.Server.CertFile != "" && a.Server.KeyFile != ""
}
