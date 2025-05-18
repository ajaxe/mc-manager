//go:build windows

package gameserver

import (
	"net/http"

	"github.com/ajaxe/mc-manager/internal/config"
	"github.com/docker/cli/cli/connhelper"
	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
)

func defaultDockerCli(logger echo.Logger) (cli *client.Client, err error) {
	cli, err = dockerCliWindows()
	logger.Info("Fetching docker client on windows")
	return
}

func dockerCliWindows() (cli *client.Client, err error) {
	c := config.LoadAppConfig()
	helper, err := connhelper.GetConnectionHelper(c.Server.DockerHostURL)

	if err != nil {
		return
	}

	httpClient := &http.Client{
		// No tls
		// No proxy
		Transport: &http.Transport{
			DialContext: helper.Dialer,
		},
	}

	var clientOpts []client.Opt

	clientOpts = append(clientOpts,
		client.WithHTTPClient(httpClient),
		client.WithHost(helper.Host),
		client.WithDialContext(helper.Dialer),
	)

	cli, err = dockerCli(clientOpts)

	return
}
