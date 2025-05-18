package gameserver

import (
	"strings"

	"github.com/ajaxe/mc-manager/internal/config"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

const LabelImageName = "mc.manager.image.name"
const LabelWordId = "mc.manager.world.id"

func defaultConfig(worldId string) container.Config {
	c := config.LoadAppConfig()

	replacer := strings.NewReplacer("${HOSTING_DIR}", c.GameServer.HostingDir)

	vols := make(map[string]struct{})
	for _, v := range c.GameServer.Volumes {
		vols[replacer.Replace(v)] = struct{}{}
	}

	labels := make(map[string]string)
	for _, l := range c.GameServer.Labels {
		splits := strings.SplitN(l, "=", 2)
		labels[splits[0]] = splits[1]
	}
	labels[LabelImageName] = c.GameServer.ImageName
	labels[LabelWordId] = worldId

	return container.Config{
		Image:   c.GameServer.ImageName,
		Volumes: vols,
		Env:     c.GameServer.EnvVars,
		Labels:  labels,
	}
}

func defaultHostConfig() container.HostConfig {
	c := config.LoadAppConfig()
	return container.HostConfig{
		LogConfig: container.LogConfig{
			Type:   c.GameServer.Logging.Driver,
			Config: c.GameServer.Logging.Options,
		},
	}
}

func defaultNetworkingConfig() network.NetworkingConfig {
	c := config.LoadAppConfig()
	networks := make(map[string]*network.EndpointSettings)
	for _, n := range c.GameServer.Networks {
		networks[n] = &network.EndpointSettings{}
	}
	return network.NetworkingConfig{
		EndpointsConfig: networks,
	}
}
