package gameserver

import (
	"context"
	"fmt"
	"strings"

	"github.com/ajaxe/mc-manager/internal/config"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func CreateGameServer(w *models.WorldItem) (resp container.CreateResponse, err error) {
	cli, err := defaultDockerCli()
	if err != nil {
		return
	}
	defer cli.Close()
	resp, err = createGameServerInternal(w, cli)
	return
}

// / return docker container name which is running the configured Image
func GameServerIntance() (name []string, err error) {
	cli, err := defaultDockerCli()
	defer cli.Close()
	if err != nil {
		return
	}

	c := config.LoadAppConfig()

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", fmt.Sprintf("%s=%s", LabelImageName, c.GameServer.ImageName))),
	})

	if err != nil {
		return
	}

	if len(containers) == 0 {
		return []string{}, nil
	}

	n := []string{}

	for _, c := range containers {
		n = append(n, strings.TrimPrefix(c.Names[0], "/"))
	}

	return n, nil
}

func createGameServerInternal(w *models.WorldItem, cli *client.Client) (resp container.CreateResponse, err error) {
	ctx := context.Background()

	c := defaultConfig()
	h := defaultHostConfig()
	n := defaultNetworkingConfig()

	resp, err = cli.ContainerCreate(ctx, &c, &h, &n, nil, ToContainerName(w.Name))

	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})

	return
}

func dockerCli(opts []client.Opt) (cli *client.Client, err error) {

	o := append(opts, client.WithAPIVersionNegotiation())

	cli, err = client.NewClientWithOpts(o...)

	return
}
