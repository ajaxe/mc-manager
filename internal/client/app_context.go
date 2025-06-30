package client

import (
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type AppContext struct {
	app.Context
}

func NewAppContext(ctx app.Context) AppContext {
	return AppContext{ctx}
}

func (c AppContext) LoadData(key string, v ...any) {
	switch key {
	case StateKeyWorlds:
		c.loadWorlds()
	case StateKeyLaunches:
		c.loadLaunches(v)
	case StateKeyCurrentPlayTimer:
		c.loadPlayTimer()
	default: // do nothing
	}
}
func (c AppContext) loadPlayTimer() {
	c.Async(func() {
		d, _ := PlayTimer()
		var s *models.PlayTimerItem
		if len(d.Data) > 0 {
			s = d.Data[0]
		} else {
			s = &models.PlayTimerItem{
				IsActive: false,
			}
		}
		c.SetState(StateKeyCurrentPlayTimer, s)
	})
}
func (c AppContext) loadWorlds() {
	c.Async(func() {
		l, _ := WorldsList()

		c.SetState(StateKeyWorlds, l.Data)
	})
}
func (c AppContext) loadLaunches(v []any) {
	req := models.LaunchItemListRequest{
		Direction: models.PageDirectionNext,
	}
	if len(v) == 1 {
		if r, ok := v[0].(models.LaunchItemListRequest); ok {
			req = r
		}
	}
	c.Async(func() {
		l, _ := LaunchList(req)

		c.SetState(StateKeyLaunches, l)
	})
}
func (c AppContext) ShowErrorMessage(r *models.ApiResult, e error) {
	if r == nil {
		r = &models.ApiResult{Success: true}
	}
	c.ShowMessage("", *r, e)
}
func (c AppContext) ShowMessage(msg string, r models.ApiResult, e error) {
	if e == nil && r.Success && msg != "" {
		c.NewActionWithValue(ActionStatusToast, StatusToastData{
			Status:  ToastStatusSuccess,
			Message: msg,
		})
	} else if e != nil || r.Success == false {
		m := r.ErrorMessage
		if e != nil {
			m = e.Error()
		}
		c.NewActionWithValue(ActionStatusToast, StatusToastData{
			Status:  ToastStatusError,
			Message: m,
		})
	} else {
		app.Logf("show message: unhandled case")
	}
	c.Update()
}
func (c AppContext) LoginRedirect(u string) {
	app.Window().
		Get("location").
		Call("replace", u)
}
