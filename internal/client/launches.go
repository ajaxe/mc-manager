package client

import (
	"github.com/ajaxe/mc-manager/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func LaunchWorld(id bson.ObjectID) (err error) {
	r := &models.ApiResult{}
	err = httpPost(buildApiURL(appBaseURL(), "/launches"), &models.CreateLaunchItem{
		WorldItemID: id,
	}, r)
	return
}
func LaunchList() (l *models.LaunchItemListResult, err error) {
	l = &models.LaunchItemListResult{}
	err = httpGet(buildApiURL(appBaseURL(), "/launches"), l)

	return
}
