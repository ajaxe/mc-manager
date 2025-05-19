package client

import (
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
func LaunchList() (l *models.LaunchItemListResult, err error) {
	l = &models.LaunchItemListResult{}
	err = httpGet(buildApiURL(appBaseURL(), "/launches"), l)

	return
}
