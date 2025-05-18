package handlers

import "github.com/ajaxe/mc-manager/internal/models"

type GameService interface {
	createGameServer(w *models.WorldItem) error
	serverIntance() ([]string, error)
	stopAllInstances() error
}
