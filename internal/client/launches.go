package client

import (
	"fmt"

	"github.com/ajaxe/mc-manager/internal/models"
)

func LaunchWorld(w *models.WorldItem) (err error) {
	r := &models.ApiResult{}
	err = httpPost(buildApiURL(appBaseURL(), "/launches"), &models.CreateLaunchItem{
		WorldItemID: w.ID,
		GameMode:    w.GameMode,
	}, r)
	return
}
func LaunchList(r models.LaunchItemListRequest) (l models.LaunchItemListResult, err error) {
	l = models.LaunchItemListResult{}
	p := fmt.Sprintf("/launches?dir=%s&cursorId=%s", r.Direction, r.CursorID)
	err = httpGet(buildApiURL(appBaseURL(), p), &l)

	return
}
