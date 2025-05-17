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
	resp, err = createGameServerInternal(w, []client.Opt{client.FromEnv})
	return
}

// / return docker container name which is running the configured Image
func GameServerIntance() (name string, err error) {
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
		return "", nil
	}

	if len(containers) > 1 {
		return "", fmt.Errorf("multiple containers found with image %s", c.GameServer.ImageName)
	}

	name = strings.TrimPrefix(containers[0].Names[0], "/")

	return
}

func createGameServerInternal(w *models.WorldItem, opts []client.Opt) (resp container.CreateResponse, err error) {
	ctx := context.Background()

	cli, err := dockerCli(opts)
	defer cli.Close()

	c := defaultConfig()
	h := defaultHostConfig()
	n := defaultNetworkingConfig()

	resp, err = cli.ContainerCreate(ctx, &c, &h, &n, nil, strings.ToLower(w.Name))

	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})

	return
}

func dockerCli(opts []client.Opt) (cli *client.Client, err error) {

	o := append(opts, client.WithAPIVersionNegotiation())

	cli, err = client.NewClientWithOpts(o...)

	return
}
