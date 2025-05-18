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

// Create creates a new game server container bades on the world item name.
func (g *GameServerOperations) Create(w *models.WorldItem) (resp container.CreateResponse, err error) {
	cli, err := defaultDockerCli(g.Logger)
	if err != nil {
		return
	}
	defer cli.Close()
	resp, err = createGameServerInternal(w, cli)
	return
}

// Intances returns List of docker container names which is running the configured Image
func (g *GameServerOperations) Intances() (name []string, err error) {
	containers, err := g.Details()
	n := []string{}
	for _, c := range containers {
		n = append(n, strings.TrimPrefix(c.Name, "/"))
	}

	return n, nil
}

func (g *GameServerOperations) Stop(w *models.WorldItem) (err error) {
	name := ToContainerName(w.Name)

	containers, err := g.Details()
	if err != nil {
		return
	}

	if len(containers) == 0 {
		return
	}

	cli, err := defaultDockerCli(g.Logger)
	if err != nil {
		return
	}
	defer cli.Close()

	for _, c := range containers {
		if strings.TrimPrefix(c.WorldID, "/") == name {
			err = cli.ContainerStop(context.Background(), c.ContainerID, container.StopOptions{})
			err = cli.ContainerRemove(context.Background(), c.ContainerID, container.RemoveOptions{Force: true})
		}
	}

	return
}

func (g *GameServerOperations) StopAll() (err error) {
	containers, err := g.Details()

	cli, err := defaultDockerCli(g.Logger)
	if err != nil {
		return
	}
	defer cli.Close()

	for _, c := range containers {
		if c.WorldID != "" {
			err = cli.ContainerStop(context.Background(), c.ContainerID, container.StopOptions{})
			err = cli.ContainerRemove(context.Background(), c.ContainerID, container.RemoveOptions{Force: true})
		}
	}

	return
}

func (g *GameServerOperations) Details() (details []*models.GameServerDetail, err error) {
	cli, err := defaultDockerCli(g.Logger)
	if err != nil {
		return
	}
	defer cli.Close()

	containers, err := listContainers(cli)
	if err != nil {
		return
	}

	if len(containers) == 0 {
		return []*models.GameServerDetail{}, nil
	}

	n := []*models.GameServerDetail{}

	for _, c := range containers {
		n = append(n, &models.GameServerDetail{
			Name:    strings.TrimPrefix(c.Names[0], "/"),
			WorldID: c.Labels[LabelWordId],
		})
	}

	return n, nil
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
