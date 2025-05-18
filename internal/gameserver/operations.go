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
	"github.com/labstack/echo/v4"
)

type GameServerOperations struct {
	Logger echo.Logger
}

// CreateGameServer creates a new game server container bades on the world item name.
func (g *GameServerOperations) CreateGameServer(w *models.WorldItem) (resp container.CreateResponse, err error) {
	cli, err := defaultDockerCli()
	if err != nil {
		return
	}
	defer cli.Close()
	resp, err = createGameServerInternal(w, cli)
	return
}

// GameServerIntance returns List of docker container names which is running the configured Image
func (g *GameServerOperations) GameServerIntance() (name []string, err error) {
	cli, err := defaultDockerCli()
	defer cli.Close()
	if err != nil {
		return
	}

	containers, err := listContainers(cli)
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

func (g *GameServerOperations) StopGameServer(w *models.WorldItem) (err error) {
	cli, err := defaultDockerCli()
	if err != nil {
		return
	}
	defer cli.Close()

	name := ToContainerName(w.Name)

	containers, err := listContainers(cli)
	if err != nil {
		return
	}

	if len(containers) == 0 {
		return
	}

	for _, c := range containers {
		if strings.TrimPrefix(c.Names[0], "/") == name {
			err = cli.ContainerStop(context.Background(), c.ID, container.StopOptions{})
			err = cli.ContainerRemove(context.Background(), c.ID, container.RemoveOptions{Force: true})
		}
	}

	return
}

func listContainers(cli *client.Client) (containers []container.Summary, err error) {
	c := config.LoadAppConfig()

	containers, err = cli.ContainerList(context.Background(), container.ListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", fmt.Sprintf("%s=%s", LabelImageName, c.GameServer.ImageName))),
	})

	return
}

func createGameServerInternal(w *models.WorldItem, cli *client.Client) (resp container.CreateResponse, err error) {
	ctx := context.Background()

	c := defaultConfig(w.ID.Hex())
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
