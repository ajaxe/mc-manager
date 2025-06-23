package gameserver

import (
	"context"

	"github.com/labstack/echo/v4"
)

type ServerCommand struct {
	Logger     echo.Logger
	GameServer *GameServerOperations
}

func (sc *ServerCommand) Message(m string) (err error) {
	cli, err := defaultDockerCli(sc.Logger)
	if err != nil {
		return
	}
	defer cli.Close()

	containers, err := sc.GameServer.listContainers(cli)

	if len(containers) == 0 {
		sc.Logger.Warnf("no active game server instances found to send message: %s", m)
	}

	for _, r := range containers {
		e := sendChatMessage(chatMessageOptions{
			gameserverConsoleOptions: gameserverConsoleOptions{
				ctx:         context.Background(),
				logger:      sc.Logger,
				cli:         cli,
				containerId: r.ID,
			},
			message: m,
		})
		if e != nil {
			sc.Logger.Errorf("failed to send message to %s: %v", r.ID, e)
		} else {
			sc.Logger.Infof("message sent to %s: %s", r.ID, m)
		}
	}

	return
}
