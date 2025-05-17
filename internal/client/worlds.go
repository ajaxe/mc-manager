package client

import (
	"github.com/ajaxe/mc-manager/internal/models"
)

func WorldsList() (d models.WorldItemListResult, err error) {
	d = models.WorldItemListResult{}
	err = httpGet(buildApiURL(appBaseURL(), "/worlds"), &d)

	return
}

func WorldCreate(w *models.WorldItem) (r models.ApiIDResult, err error) {
	r = models.ApiIDResult{}
	err = httpPost(buildApiURL(appBaseURL(), "/worlds"), w, &r)
	return
}
