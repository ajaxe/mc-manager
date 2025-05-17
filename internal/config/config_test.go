package config

import "testing"

func TestReadConfig_Ok(t *testing.T) {
	config, err := loadAppConfigInternal("../../", "config")
	if err != nil {
		t.Fatalf("failed to get app config: %v", err)
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
}
