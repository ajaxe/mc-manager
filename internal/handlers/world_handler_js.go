//go:build wasm

package handlers

import (
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/labstack/echo/v4"
)

type gameService struct {
}

func NewGameService(logger echo.Logger) GameService {
	return &gameService{}
}

func (g *gameService) serverDetails() ([]*models.GameServerDetail, error) {
	return []*models.GameServerDetail{}, nil
}
func (g *gameService) createGameServer(w *models.WorldItem) (err error) {
	return nil
}
func (g *gameService) stopAllInstances() error {
	return nil
}
func toContainerName(s string) string {
	return s
}
