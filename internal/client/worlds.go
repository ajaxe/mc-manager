package client

import (
	"github.com/ajaxe/mc-manager/internal/models"
)

func WorldsList() (d models.WorldItemListResult, err error) {
	d = models.WorldItemListResult{}
	err = httpGet(buildApiURL(appBaseURL(), "/worlds"), &d)

	return
}
