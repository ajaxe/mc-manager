//go:build unix || windows

package handlers

import (
	"github.com/ajaxe/mc-manager/internal/gameserver"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/labstack/echo/v4"
)

type gameService struct {
	op *gameserver.GameServerOperations
}

func NewGameService(logger echo.Logger) GameService {
	return &gameService{
		op: &gameserver.GameServerOperations{
			Logger: logger,
		},
	}
}
func (g *gameService) gameServerIntance() (n []string, err error) {
	n, err = g.op.GameServerIntance()
	return
}

func (g *gameService) createGameServer(w *models.WorldItem) (err error) {
	_, err = g.op.CreateGameServer(w)
	return
}
func toContainerName(s string) string {
	return gameserver.ToContainerName(s)
}
