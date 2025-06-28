package config

import (
	"net/url"
)

type AppConfig struct {
	Server struct {
		// port on which the server listens for incoming connections.
		Port string `mapstructure:"port"`

		// public certificate to setup TLS for development server.
		CertFile string `mapstructure:"cert_file"`

		// private certificate to setup TLS for development server.
		KeyFile string `mapstructure:"key_file"`

		// dockerhost URL when server running on windows.
		DockerHostURL string `mapstructure:"docker_host_url"`

		// Authentication redirect host URL when hosted behind reverse proxy.
		AuthServerURL string `mapstructure:"auth_server_url"`

		// Authentication redirect URL when hosted behind reverse proxy.
		AuthRedirectPath string `mapstructure:"auth_redirect_path"`

		// Auth cookie name to check for authenticated session.
		AuthCookieName string `mapstructure:"auth_cookie_name"`

		// Token used to identify service during authentication
		AuthToken string `mapstructure:"auth_token"`
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

func (a *AppConfig) UseTLS() bool {
	return a.Server.CertFile != "" && a.Server.KeyFile != ""
}

// AuthRedirectURL combines config.Server.AuthRedirectURL & config.Server.AuthToken
// to returns authentication redirect URL.
func (a *AppConfig) AuthRedirectURL() (u string, err error) {
	if a.Server.AuthServerURL == "" || a.Server.AuthRedirectPath == "" || a.Server.AuthToken == "" {
		u = ""
		return
	}
	p, err := url.Parse(a.Server.AuthServerURL + a.Server.AuthRedirectPath)
	if err != nil {
		return
	}
	qs := p.Query()
	qs.Set("token", a.Server.AuthToken)

	p.RawQuery = qs.Encode()

	u = p.String()
	return
}
