package gameserver

import (
	"fmt"
	"strings"

	"github.com/ajaxe/mc-manager/internal/config"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/labstack/echo/v4"
)

const LabelImageName = "mc.manager.image.name"
const LabelWorldId = "mc.manager.world.id"
const EnvVarGamemode = "MC_GAMEMODE"
const EnvVarLevelName = "MC_LEVEL_NAME"
const EnvVarLevelSeed = "MC_LEVEL_SEED"

type ServiceConfig struct {
	Logger echo.Logger
	Config config.AppConfig
}

func (s *ServiceConfig) defaultConfig(w *models.WorldItem) container.Config {
	replacer := strings.NewReplacer("${HOSTING_DIR}", s.Config.GameServer.HostingDir)

	vols := make(map[string]struct{})
	for _, v := range s.Config.GameServer.Volumes {
		vols[replacer.Replace(v)] = struct{}{}
	}

	labels := make(map[string]string)
	for _, l := range s.Config.GameServer.Labels {
		splits := strings.SplitN(l, "=", 2)
		labels[splits[0]] = splits[1]
	}
	labels[LabelImageName] = s.Config.GameServer.ImageName
	labels[LabelWorldId] = w.ID.Hex()

	env := []string{}
	for _, v := range s.Config.GameServer.EnvVars {
		if strings.HasPrefix(v, EnvVarGamemode) {
			env = append(env, fmt.Sprintf("%s=%s", EnvVarGamemode, w.GameMode))
		} else if strings.HasPrefix(v, EnvVarLevelName) {
			env = append(env, fmt.Sprintf("%s=%s", EnvVarLevelName, w.Name))
		} else if strings.HasPrefix(v, EnvVarLevelSeed) {
			env = append(env, fmt.Sprintf("%s=%s", EnvVarLevelSeed, w.WorldSeed))
		} else {
			env = append(env, v)
		}
	}

	return container.Config{
		Image:   s.Config.GameServer.ImageName,
		Volumes: vols,
		Env:     env,
		Labels:  labels,
	}
}

func (s *ServiceConfig) defaultHostConfig() container.HostConfig {
	return container.HostConfig{
		LogConfig: container.LogConfig{
			Type:   s.Config.GameServer.Logging.Driver,
			Config: s.Config.GameServer.Logging.Options,
		},
	}
}

func (s *ServiceConfig) defaultNetworkingConfig() network.NetworkingConfig {
	networks := make(map[string]*network.EndpointSettings)
	for _, n := range s.Config.GameServer.Networks {
		networks[n] = &network.EndpointSettings{}
	}
	return network.NetworkingConfig{
		EndpointsConfig: networks,
	}
}
