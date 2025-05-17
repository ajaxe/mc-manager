package config

type AppConfig struct {
	Server struct {
		Port     string `mapstructure:"port"`
		CertFile string `mapstructure:"cert_file"`
		KeyFile  string `mapstructure:"key_file"`
	} `mapstructure:"server"`
	Database struct {
		ConnectionURI string `mapstructure:"connection_uri"`
		DbName        string `mapstructure:"db_name"`
	} `mapstructure:"database"`
}

func (a AppConfig) UseTLS() bool {
	return a.Server.CertFile != "" && a.Server.KeyFile != ""
}
