//go:build unix || windows

package handlers

import "github.com/ajaxe/mc-manager/internal/gameserver"

func gameServerIntance() (n string, err error) {
	n, err = gameserver.GameServerIntance()
	return
}
