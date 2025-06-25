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

func NewGameServerOperations(l echo.Logger, cfg config.AppConfig) *GameServerOperations {
	return &GameServerOperations{
		Logger: l,
		Config: &ServiceConfig{
			Logger: l,
			Config: cfg,
		},
	}
}

type GameServerOperations struct {
	Logger echo.Logger
	Config *ServiceConfig
}

// Create creates a new game server container bades on the world item name.
func (g *GameServerOperations) Create(w *models.WorldItem) (resp container.CreateResponse, err error) {
	cli, err := defaultDockerCli(g.Logger)
	if err != nil {
		return
	}
	defer cli.Close()
	resp, err = g.createGameServerInternal(w, cli)
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
			g.removeContainer(c, cli)
		}
	}

	return
}

// StopAll stops & removes all containers running the configured Image with World ID label.
func (g *GameServerOperations) StopAll() (err error) {
	containers, err := g.Details()

	cli, err := defaultDockerCli(g.Logger)
	if err != nil {
		return
	}
	defer cli.Close()

	for _, c := range containers {
		if c.WorldID != "" {
			g.removeContainer(c, cli)
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

	containers, err := g.listContainers(cli)
	if err != nil {
		return
	}

	if len(containers) == 0 {
		return []*models.GameServerDetail{}, nil
	}

	n := []*models.GameServerDetail{}

	for _, c := range containers {
		n = append(n, &models.GameServerDetail{
			Name:        strings.TrimPrefix(c.Names[0], "/"),
			WorldID:     c.Labels[LabelWorldId],
			GameMode:    c.Labels[LabelWorldGameMode],
			ContainerID: c.ID,
		})
	}

	return n, nil
}

func (g *GameServerOperations) listContainers(cli *client.Client) (containers []container.Summary, err error) {
	c := g.Config.Config
	l := fmt.Sprintf("%s=%s", LabelImageName, c.GameServer.ImageName)

	g.Logger.Infof("listing container with label: %v", l)

	containers, err = cli.ContainerList(context.Background(), container.ListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", l)),
	})

	return
}

func (g *GameServerOperations) createGameServerInternal(w *models.WorldItem, cli *client.Client) (resp container.CreateResponse, err error) {
	ctx := context.Background()

	c := g.Config.defaultConfig(w)
	h := g.Config.defaultHostConfig()
	n := g.Config.defaultNetworkingConfig()

	cname := ToContainerName(w.Name)

	resp, err = cli.ContainerCreate(ctx, &c, &h, &n, nil, cname)
	if err != nil {
		return
	}

	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})

	if err == nil {
		g.Logger.Infof("container started, name:%s id:%s", cname, resp.ID)
	}

	return
}

func (g *GameServerOperations) removeContainer(c *models.GameServerDetail, cli *client.Client) (err error) {
	g.Logger.Infof("stopping container:%s", c.ContainerID)
	err = sendStopCommand(gameserverConsoleOptions{
		ctx:         context.Background(),
		logger:      g.Logger,
		cli:         cli,
		containerId: c.ContainerID,
	})

	if err != nil {
		return err
	} else {
		g.Logger.Infof("removing container:%s", c.ContainerID)
		err = cli.ContainerRemove(context.Background(), c.ContainerID, container.RemoveOptions{Force: true})
	}
	return
}

func dockerCli(opts []client.Opt) (cli *client.Client, err error) {

	o := append(opts, client.WithAPIVersionNegotiation())

	cli, err = client.NewClientWithOpts(o...)

	return
}

func (g *GameServerOperations) Message(m string) (err error) {
	cli, err := defaultDockerCli(g.Logger)
	if err != nil {
		return
	}
	defer cli.Close()

	containers, err := g.listContainers(cli)

	if len(containers) == 0 {
		g.Logger.Warnf("no active game server instances found to send message: %s", m)
	}

	for _, r := range containers {
		e := sendChatMessage(chatMessageOptions{
			gameserverConsoleOptions: gameserverConsoleOptions{
				ctx:         context.Background(),
				logger:      g.Logger,
				cli:         cli,
				containerId: r.ID,
			},
			message: m,
		})
		if e != nil {
			g.Logger.Errorf("failed to send message to %s: %v", r.ID, e)
		} else {
			g.Logger.Infof("message sent to %s: %s", r.ID, m)
		}
	}

	return
}