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
	return AppConfig{
		Server: Server{
			Port:     "8000",
			CertFile: "F:\\mkcert\\certificates\\localhost+2.pem",
			KeyFile:  "F:\\mkcert\\certificates\\localhost+2-key.pem",
		},
	}
}

func (a AppConfig) UseTLS() bool {
	return a.Server.CertFile != "" && a.Server.KeyFile != ""
}
