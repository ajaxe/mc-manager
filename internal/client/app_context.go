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
	case StateKeyLaunches:
		c.loadLaunches()
	default: // do nothing
	}
}

func (c AppContext) loadWorlds() {
	c.Async(func() {
		l, _ := WorldsList()

		c.SetState(StateKeyWorlds, l.Data)
	})
}
func (c AppContext) loadLaunches() {
	c.Async(func() {
		l, _ := LaunchList()

		c.SetState(StateKeyLaunches, l.Data)
	})
}
