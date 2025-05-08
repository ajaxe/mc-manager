package models

type AppConfig struct {
	Server Server `mapstructure:"server"`
}
type Server struct {
	Port     string `mapstructure:"port"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

func LoadAppConfig() AppConfig {
	return AppConfig{}
}

func (a AppConfig) UseTLS() bool {
	return a.Server.CertFile != "" && a.Server.KeyFile != ""
}
