package gameserver

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
)

type ServerCommand struct {
	Logger echo.Logger
}

func (sc *ServerCommand) Message(m string) error {
	return nil
}

func defaultAttachOptions() container.AttachOptions {
	return container.AttachOptions{
		Stream: false,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Logs:   false,
	}
}

// sendChatMessage executes "say" command to send a message to all playes on the game server.
func sendChatMessage(m string, ctx context.Context, containerId string, cli *client.Client) (err error) {
	r, err := cli.ContainerAttach(ctx, containerId, defaultAttachOptions())
	if err != nil {
		return
	}
	defer r.Close()

	_, err = r.Conn.Write([]byte(fmt.Sprintf("say %s", m)))

	return err
}

// sendStopCommand executes "stop" command to gracefully stop the game server.
func sendStopCommand(ctx context.Context, containerId string, cli *client.Client) (err error) {
	r, err := cli.ContainerAttach(ctx, containerId, defaultAttachOptions())
	if err != nil {
		return
	}
	defer r.Close()

	_, err = r.Conn.Write([]byte("stop"))

	return err
}
