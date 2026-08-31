//go:build wasm

package handlers

import (
	"log/slog"

	"github.com/ajaxe/mc-manager/internal/models"
)

type gameService struct {
}

func NewGameService(logger *slog.Logger) GameService {
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
func (g *gameService) sendMessageToServer(message string) (err error) {
	return nil
}
