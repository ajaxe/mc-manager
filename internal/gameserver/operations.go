package gameserver

import (
	"context"
	"strings"

	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func CreateGameServer(w *models.WorldItem) (resp container.CreateResponse, err error) {
	resp, err = createGameServerInternal(w, []client.Opt{client.FromEnv})
	return
}
func createGameServerInternal(w *models.WorldItem, opts []client.Opt) (resp container.CreateResponse, err error) {
	ctx := context.Background()

	o := append(opts, client.WithAPIVersionNegotiation())

	cli, err := client.NewClientWithOpts(o...)
	if err != nil {
		return
	}
	defer cli.Close()

	c := defaultConfig()
	h := defaultHostConfig()
	n := defaultNetworkingConfig()

	resp, err = cli.ContainerCreate(ctx, &c, &h, &n, nil, strings.ToLower(w.Name))

	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})

	return
}
