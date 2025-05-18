//go:build unix

package gameserver

import (
	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
)

func defaultDockerCli(logger echo.Logger) (cli *client.Client, err error) {
	cli, err = dockerCli([]client.Opt{client.FromEnv})
	logger.Info("Fetching docker client on linux")
	return
}
