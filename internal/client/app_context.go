package client

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type AppContext struct {
	app.Context
}

func NewAppContext(ctx app.Context) AppContext {
	return AppContext{ctx}
}

func (c AppContext) LoadData(key string) {
	switch key {
	case StateKeyWorlds:
		c.loadWorlds()
	default: // do nothing
	}
}

func (c AppContext) loadWorlds() {
	c.Async(func() {
		l, _ := WorldsList()

		for _, item := range l.Data {
			app.Logf("load worlds %v", item)
		}
		c.SetState(StateKeyWorlds, l.Data)
	})
}
