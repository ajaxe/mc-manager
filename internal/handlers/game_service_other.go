//go:build unix || windows

package handlers

import (
	"github.com/ajaxe/mc-manager/internal/config"
	"github.com/ajaxe/mc-manager/internal/gameserver"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/labstack/echo/v4"
)

type gameService struct {
	op *gameserver.GameServerOperations
}

func NewGameService(logger echo.Logger) GameService {
	ops := gameserver.NewGameServerOperations(logger, config.LoadAppConfig())

	return &gameService{
		op: ops,
	}
}
func (g *gameService) serverDetails() (n []*models.GameServerDetail, err error) {
	n, err = g.op.Details()
	return
}

func (g *gameService) createGameServer(w *models.WorldItem) (err error) {
	_, err = g.op.Create(w)
	return
}
func (g *gameService) stopAllInstances() error {
	return g.op.StopAll()
}
func (g *gameService) sendMessageToServer(message string) (err error) {
	err = g.op.Message(message)
	return
}

func toContainerName(s string) string {
	return gameserver.ToContainerName(s)
}
