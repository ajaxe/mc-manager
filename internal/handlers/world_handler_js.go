//go:build js

package handlers

import "github.com/ajaxe/mc-manager/internal/models"

func gameServerIntance() ([]string, error) {
	return []string{}, nil
}
func createGameServer(w *models.WorldItem) (err error) {
	return nil
}
func toContainerName(s string) string {
	return s
}
