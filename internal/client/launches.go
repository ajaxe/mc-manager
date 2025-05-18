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
