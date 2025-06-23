package gameserver

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
)

type gameserverConsoleOptions struct {
	ctx         context.Context
	logger      echo.Logger
	cli         *client.Client
	containerId string
}

type chatMessageOptions struct {
	gameserverConsoleOptions
	message string
}

func defaultAttachOptions() container.AttachOptions {
	return container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Logs:   false,
	}
}

// sendChatMessage executes "say" command to send a message to all playes on the game server.
func sendChatMessage(opts chatMessageOptions) (err error) {
	r, err := opts.cli.ContainerAttach(opts.ctx, opts.containerId, defaultAttachOptions())
	if err != nil {
		return
	}
	defer r.Close()
	mt, ok := r.MediaType()
	opts.logger.Debugf("hijacked response media-type: %s: ok-%v", mt, ok)
	m := fmt.Sprintf("say %s\n", opts.message)
	opts.logger.Debugf("sending message to %s: %s", opts.containerId, m)

	_, err = r.Conn.Write([]byte(m))

	return err
}

// sendStopCommand executes "stop" command to gracefully stop the game server.
func sendStopCommand(opts gameserverConsoleOptions) (err error) {
	r, err := opts.cli.ContainerAttach(opts.ctx, opts.containerId, defaultAttachOptions())
	if err != nil {
		return
	}
	defer r.Close()

	_, err = r.Conn.Write([]byte("stop"))

	return err
}
