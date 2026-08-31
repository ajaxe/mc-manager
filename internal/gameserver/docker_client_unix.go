//go:build unix

package gameserver

import (
	"log/slog"

	"github.com/docker/docker/client"
)

func defaultDockerCli(logger *slog.Logger) (cli *client.Client, err error) {
	cli, err = dockerCli([]client.Opt{client.FromEnv})
	logger.Info("Fetching docker client on linux")
	return
}
