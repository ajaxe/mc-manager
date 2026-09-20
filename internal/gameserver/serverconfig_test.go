package gameserver

import (
	"fmt"
	"io"
	"log/slog"
	"slices"
	"testing"

	"github.com/ajaxe/mc-manager/internal/config"
	"github.com/ajaxe/mc-manager/internal/models"
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

func TestDefaultHostConfigResources(t *testing.T) {
	sut := createServerConfig()

	// Override config with known test values
	sut.Config.GameServer.Deploy.Resources.Limits.Memory = "512m"
	sut.Config.GameServer.Deploy.Resources.Limits.CPUs = "0.50"
	sut.Config.GameServer.Deploy.Resources.Reservations.Memory = "128m"
	sut.Config.GameServer.Deploy.Resources.Reservations.CPUs = "0.10"

	hostConfig := sut.defaultHostConfig()

	expectedMemory := int64(512 * 1024 * 1024)
	if hostConfig.Resources.Memory != expectedMemory {
		t.Errorf("expected Memory %d, got: %d", expectedMemory, hostConfig.Resources.Memory)
	}

	expectedNanoCPUs := int64(500000000)
	if hostConfig.Resources.NanoCPUs != expectedNanoCPUs {
		t.Errorf("expected NanoCPUs %d, got: %d", expectedNanoCPUs, hostConfig.Resources.NanoCPUs)
	}

	expectedMemoryReservation := int64(128 * 1024 * 1024)
	if hostConfig.Resources.MemoryReservation != expectedMemoryReservation {
		t.Errorf("expected MemoryReservation %d, got: %d", expectedMemoryReservation, hostConfig.Resources.MemoryReservation)
	}

	expectedCPUShares := int64(102)
	if hostConfig.Resources.CPUShares != expectedCPUShares {
		t.Errorf("expected CPUShares %d, got: %d", expectedCPUShares, hostConfig.Resources.CPUShares)
	}
}

func createServerConfig() ServiceConfig {
	l := slog.New(slog.NewJSONHandler(io.Discard, nil))
	return ServiceConfig{
		Logger: l,
		Config: config.LoadAppConfig(),
	}
}
