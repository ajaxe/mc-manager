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
	return &gameService{
		op: &gameserver.GameServerOperations{
			Logger: logger,
			Config: &gameserver.ServiceConfig{
				Logger: logger,
				Config: config.LoadAppConfig(),
			},
		},
	}
}
func (g *gameService) serverIntance() (n []string, err error) {
	n, err = g.op.Intances()
	return
}

func (g *gameService) createGameServer(w *models.WorldItem) (err error) {
	_, err = g.op.Create(w)
	return
}
func (g *gameService) stopAllInstances() error {
	return g.op.StopAll()
}

func toContainerName(s string) string {
	return gameserver.ToContainerName(s)
}
