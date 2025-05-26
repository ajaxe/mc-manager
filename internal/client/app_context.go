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
	default: // do nothing
	}
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
