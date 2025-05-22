package gameserver

import (
	"fmt"
	"slices"
	"testing"

	"github.com/ajaxe/mc-manager/internal/config"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestDefaultConfig(t *testing.T) {
	n := newWorldName("TestDefaultConfig")
	seed := "10000"
	gamemode := "survival"

	data := &models.WorldItem{
		ID:        bson.NewObjectID().Hex(),
		Name:      n,
		WorldSeed: seed,
		GameMode:  gamemode,
	}

	sut := createServerConfig()

	cfg := sut.defaultConfig(data)

	if cfg.Labels[LabelWorldId] != data.ID {
		t.Fatalf("Label '%s' must be '%s', got: '%v'", LabelWorldId, data.ID, cfg.Labels[LabelWorldId])
	}
	if cfg.Labels[LabelImageName] != sut.Config.GameServer.ImageName {
		t.Fatalf("Label '%s' must be '%s', got: '%v'", LabelImageName,
			sut.Config.GameServer.ImageName, cfg.Labels[LabelImageName])
	}

	t.Logf("Env vars: %v", cfg.Env)

	if !slices.Contains(cfg.Env, fmt.Sprintf("%s=%s", EnvVarGamemode, data.GameMode)) {
		t.Fatalf("EnvVar '%s' must be '%s', got: '%v'", EnvVarGamemode,
			data.GameMode, cfg.Env)
	}
	if !slices.Contains(cfg.Env, fmt.Sprintf("%s=%s", EnvVarLevelName, data.Name)) {
		t.Fatalf("EnvVar '%s' must be '%s', got: '%v'", EnvVarLevelName,
			data.Name, cfg.Env)
	}
	if !slices.Contains(cfg.Env, fmt.Sprintf("%s=%s", EnvVarLevelSeed, data.WorldSeed)) {
		t.Fatalf("EnvVar '%s' must be '%s', got: '%v'", EnvVarLevelSeed,
			data.WorldSeed, cfg.Env)
	}
}

func createServerConfig() ServiceConfig {
	l := log.New("echo_test")
	return ServiceConfig{
		Logger: l,
		Config: config.LoadAppConfig(),
	}
}
