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

func WorldDelete(id string) (r models.ApiResult, err error) {
	r = models.ApiResult{}
	err = httpDelete(buildApiURL(appBaseURL(), "/worlds/"+id), &r)
	return
}
func WorldUpdate(w *models.WorldItem) (r models.ApiResult, err error) {
	r = models.ApiResult{}
	err = httpPut(buildApiURL(appBaseURL(), "/worlds/"+w.ID), w, &r)
	return
}
