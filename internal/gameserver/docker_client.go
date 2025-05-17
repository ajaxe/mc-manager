// +build windows

package gameserver

import (
	"net/http"

	"github.com/docker/cli/cli/connhelper"
	"github.com/docker/docker/client"
)

func defaultDockerCli() (cli *client.Client, err error) {
	cli, err = dockerCliWindows()
	return
}

func dockerCliWindows() (cli *client.Client, err error) {
	helper, err := connhelper.GetConnectionHelper("ssh://ajay@dockerhost.local")

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
