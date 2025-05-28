package config

import (
	"fmt"
	"os"
	"testing"
)

func TestReadConfig_Ok(t *testing.T) {
	token := "1234"
	os.Setenv("APP_SERVER_AUTH_TOKEN", token)
	config, err := loadAppConfigInternal("../../", "config")
	if err != nil {
		t.Fatalf("failed to get app config: %v", err)
	}
	if config.Server.AuthToken != token {
		t.Fatalf("Invalid auth token. expected: %s, got: [%s]", token, config.Server.AuthToken)
	}
	if config.Server.Port == "" {
		t.Fatal("Server port is empty")
	}
	if config.Database.ConnectionURI == "" {
		t.Fatal("Database connection URI is empty")
	}
	if config.GameServer.HostingDir == "" {
		t.Fatal("Game server hosting dir is empty")
	}
	if config.GameServer.ImageName == "" {
		t.Fatal("Game server image name is empty")
	}
	if config.GameServer.EnvVars == nil {
		t.Fatal("Game server env vars is nil")
	}
	if len(config.GameServer.EnvVars) == 0 {
		t.Fatal("Game server env vars is empty")
	}
	if config.GameServer.Volumes == nil {
		t.Fatal("Game server volumes is nil")
	}
	if len(config.GameServer.Volumes) == 0 {
		t.Fatal("Game server volumes is empty")
	}
	if config.GameServer.Labels == nil {
		t.Fatal("Game server labels is nil")
	}
	if len(config.GameServer.Labels) == 0 {
		t.Fatal("Game server labels is empty")
	}
	if config.GameServer.Networks == nil {
		t.Fatal("Game server networks is nil")
	}
	if len(config.GameServer.Networks) == 0 {
		t.Fatal("Game server networks is empty")
	}
	if config.GameServer.Logging.Driver == "" {
		t.Fatal("Game server logging driver is empty")
	}
	if config.GameServer.Logging.Options == nil {
		t.Fatal("Game server logging options is nil")
	}
	if config.GameServer.Logging.Options["max-size"] == "" {
		t.Fatal("Game server logging max-size option is empty")
	}
	if config.GameServer.Logging.Options["loki-url"] == "" {
		t.Fatal("Game server logging loki-url option is empty")
	}
	os.Unsetenv("APP_SERVER_AUTH_TOKEN")
}

func TestAuthRedirectURL(t *testing.T) {
	token := "1234"
	url := "https://localhost-test.com"

	os.Setenv("APP_SERVER_AUTH_TOKEN", token)
	os.Setenv("APP_SERVER_AUTH_REDIRECT_URL", url)

	config, err := loadAppConfigInternal("../../", "config")

	if err != nil {
		t.Fatalf("error loading config: %v", err)
	}

	expected := fmt.Sprintf("%s?token=%s", url, token)
	u, err := config.AuthRedirectURL()

	if err != nil {
		t.Fatalf("error parsing auth config: %v", err)
	}

	if u != expected {
		t.Fatalf("invlaid auth redirect url. expected: %s, got: [%s]", expected, u)
	}
}
