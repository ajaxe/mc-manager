//go:build js

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

func (g *gameService) serverIntance() ([]string, error) {
	return []string{}, nil
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
