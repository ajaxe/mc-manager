package pages

import (
	"github.com/ajaxe/mc-manager/internal/client"
	"github.com/ajaxe/mc-manager/internal/components"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type PlayTimerPage struct {
	app.Compo
}

func (h *PlayTimerPage) OnNav(ctx app.Context) {
	client.NewAppContext(ctx).
		LoadData(client.StateKeyCurrentPlayTimer)
}
func (p *PlayTimerPage) Render() app.UI {
	return &MainLayout{
		Content: []app.UI{
			components.AppPlayTimer(),
		},
	}
}
