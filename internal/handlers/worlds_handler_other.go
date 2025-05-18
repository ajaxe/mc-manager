//go:build unix || windows

package handlers

import (
	"github.com/ajaxe/mc-manager/internal/gameserver"
	"github.com/ajaxe/mc-manager/internal/models"
)

func gameServerIntance() (n []string, err error) {
	n, err = gameserver.GameServerIntance()
	return
}

func createGameServer(w *models.WorldItem) (err error) {
	_, err = gameserver.CreateGameServer(w)
	return
}
func toContainerName(s string) string {
	return gameserver.ToContainerName(s)
}
