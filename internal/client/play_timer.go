package client

import (
	"github.com/ajaxe/mc-manager/internal/models"
)

func PlayTimer() (d models.PlayTimerListResult, err error) {
	d = models.PlayTimerListResult{}
	err = httpGet(buildApiURL(appBaseURL(), "/playtimer"), &d)

	return
}

func StopPlaytimer() (err error) {
	err = httpDelete(buildApiURL(appBaseURL(), "/playtimer"), nil)
	return
}
func StartPlaytimer(p *models.PlayTimerItem) (err error) {
	r := &models.ApiIDResult{}
	err = httpPost(buildApiURL(appBaseURL(), "/playtimer"), p, r)
	return
}
