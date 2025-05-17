// +build unix

package gameserver

import "github.com/docker/docker/client"

func defaultDockerCli() (cli *client.Client, err error) {
	cli, err = dockerCli([]client.Opt{client.FromEnv})
	return
}
