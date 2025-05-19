package handlers

import "github.com/ajaxe/mc-manager/internal/models"

type GameService interface {
	createGameServer(w *models.WorldItem) error
	serverDetails() ([]*models.GameServerDetail, error)
	stopAllInstances() error
}
