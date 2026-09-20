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
	if config.GameServer.Deploy.Resources.Limits.Memory != "512m" {
		t.Fatalf("expected deploy limits memory '512m', got: %s", config.GameServer.Deploy.Resources.Limits.Memory)
	}
	if config.GameServer.Deploy.Resources.Limits.CPUs != "0.50" {
		t.Fatalf("expected deploy limits cpus '0.50', got: %s", config.GameServer.Deploy.Resources.Limits.CPUs)
	}
	if config.GameServer.Deploy.Resources.Reservations.Memory != "128m" {
		t.Fatalf("expected deploy reservations memory '128m', got: %s", config.GameServer.Deploy.Resources.Reservations.Memory)
	}
	if config.GameServer.Deploy.Resources.Reservations.CPUs != "0.10" {
		t.Fatalf("expected deploy reservations cpus '0.10', got: %s", config.GameServer.Deploy.Resources.Reservations.CPUs)
	}
	os.Unsetenv("APP_SERVER_AUTH_TOKEN")
}

func TestAuthRedirectURL(t *testing.T) {
	token := "1234"
	url := "https://localhost-test.com"
	path := "/"
	ck := ".test.auth.cookie"

	os.Setenv("APP_SERVER_AUTH_SERVER_ENABLED", "true")
	os.Setenv("APP_SERVER_AUTH_TOKEN", token)
	os.Setenv("APP_SERVER_AUTH_SERVER_URL", url)
	os.Setenv("APP_SERVER_AUTH_REDIRECT_PATH", path)
	os.Setenv("APP_SERVER_AUTH_COOKIE_NAME", ck)

	config, err := loadAppConfigInternal("../../", "config")

	if err != nil {
		t.Fatalf("error loading config: %v", err)
	}

	expected := fmt.Sprintf("%s%s?token=%s", url, path, token)
	u, err := config.AuthRedirectURL()

	if err != nil {
		t.Fatalf("error parsing auth config: %v", err)
	}

	if u != expected {
		t.Fatalf("invalid auth redirect url. expected: %s, got: [%s]", expected, u)
	}

	os.Unsetenv("APP_SERVER_AUTH_SERVER_ENABLED")
	os.Unsetenv("APP_SERVER_AUTH_TOKEN")
	os.Unsetenv("APP_SERVER_AUTH_SERVER_URL")
	os.Unsetenv("APP_SERVER_AUTH_REDIRECT_PATH")
	os.Unsetenv("APP_SERVER_AUTH_COOKIE_NAME")
}

func TestAuthIntrospectURL(t *testing.T) {
	url := "https://localhost-test.com"
	path := "/introspect"

	os.Setenv("APP_SERVER_AUTH_SERVER_ENABLED", "true")
	os.Setenv("APP_SERVER_AUTH_SERVER_URL", url)
	os.Setenv("APP_SERVER_AUTH_INTROSPECT_PATH", path)

	config, err := loadAppConfigInternal("../../", "config")

	if err != nil {
		t.Fatalf("error loading config: %v", err)
	}

	expected := fmt.Sprintf("%s%s", url, path)
	u, err := config.AuthIntrospectURL()

	if err != nil {
		t.Fatalf("error parsing auth config: %v", err)
	}

	if u != expected {
		t.Fatalf("invalid auth introspect url. expected: %s, got: [%s]", expected, u)
	}

	os.Unsetenv("APP_SERVER_AUTH_SERVER_ENABLED")
	os.Unsetenv("APP_SERVER_AUTH_SERVER_URL")
	os.Unsetenv("APP_SERVER_AUTH_INTROSPECT_PATH")
}

func TestEnvOverride_Ok(t *testing.T) {
	ck := ".test.auth.cookie"
	t.Setenv("APP_SERVER_PORT", "9090")
	t.Setenv("APP_IS_DEV", "false")
	t.Setenv("APP_DATABASE_CONNECTION_URI", "mongodb://custom:27017")
	t.Setenv("APP_DATABASE_DB_NAME", "custom_db")
	t.Setenv("APP_GAME_SERVER_HOSTING_DIR", "/tmp/custom")
	t.Setenv("APP_GAME_SERVER_LOGGING_DRIVER", "json-file")
	t.Setenv("APP_GAME_SERVER_DEPLOY_RESOURCES_LIMITS_MEMORY", "1g")
	t.Setenv("APP_GAME_SERVER_DEPLOY_RESOURCES_LIMITS_CPUS", "1.5")
	t.Setenv("APP_GAME_SERVER_DEPLOY_RESOURCES_RESERVATIONS_MEMORY", "256m")
	t.Setenv("APP_GAME_SERVER_DEPLOY_RESOURCES_RESERVATIONS_CPUS", "0.25")
	t.Setenv("APP_SERVER_AUTH_COOKIE_NAME", ck)

	config, err := loadAppConfigInternal("../../", "config")
	if err != nil {
		t.Fatalf("failed to load app config: %v", err)
	}

	if config.Server.Port != "9090" {
		t.Errorf("expected server port 9090, got: %s", config.Server.Port)
	}
	if config.IsDev != false {
		t.Errorf("expected is_dev false, got: %v", config.IsDev)
	}
	if config.Database.ConnectionURI != "mongodb://custom:27017" {
		t.Errorf("expected connection_uri 'mongodb://custom:27017', got: %s", config.Database.ConnectionURI)
	}
	if config.Database.DbName != "custom_db" {
		t.Errorf("expected db_name 'custom_db', got: %s", config.Database.DbName)
	}
	if config.GameServer.HostingDir != "/tmp/custom" {
		t.Errorf("expected hosting_dir '/tmp/custom', got: %s", config.GameServer.HostingDir)
	}
	if config.GameServer.Logging.Driver != "json-file" {
		t.Errorf("expected logging driver 'json-file', got: %s", config.GameServer.Logging.Driver)
	}
	if config.GameServer.Deploy.Resources.Limits.Memory != "1g" {
		t.Errorf("expected deploy limits memory '1g', got: %s", config.GameServer.Deploy.Resources.Limits.Memory)
	}
	if config.GameServer.Deploy.Resources.Limits.CPUs != "1.5" {
		t.Errorf("expected deploy limits cpus '1.5', got: %s", config.GameServer.Deploy.Resources.Limits.CPUs)
	}
	if config.GameServer.Deploy.Resources.Reservations.Memory != "256m" {
		t.Errorf("expected deploy reservations memory '256m', got: %s", config.GameServer.Deploy.Resources.Reservations.Memory)
	}
	if config.GameServer.Deploy.Resources.Reservations.CPUs != "0.25" {
		t.Errorf("expected deploy reservations cpus '0.25', got: %s", config.GameServer.Deploy.Resources.Reservations.CPUs)
	}
	if config.Server.AuthCookieName != ck {
		t.Errorf("expected auth cookie name '%s', got: %s", ck, config.Server.AuthCookieName)
	}
}

func TestMissingConfigFile_EnvOnly(t *testing.T) {
	t.Setenv("APP_SERVER_PORT", "7070")
	t.Setenv("APP_DATABASE_DB_NAME", "env_db")

	config, err := loadAppConfigInternal("non_existent_dir", "non_existent_config")
	if err != nil {
		t.Fatalf("expected no error for missing config file with env vars set, got: %v", err)
	}

	if config.Server.Port != "7070" {
		t.Errorf("expected server port '7070', got: '%s'", config.Server.Port)
	}
	if config.Database.DbName != "env_db" {
		t.Errorf("expected db_name 'env_db', got: '%s'", config.Database.DbName)
	}
}
