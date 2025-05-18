package pages

import (
	"github.com/ajaxe/mc-manager/internal/client"
	"github.com/ajaxe/mc-manager/internal/components"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type WorldLaunchPage struct {
	app.Compo
}

func (h *WorldLaunchPage) OnNav(ctx app.Context) {
	client.NewAppContext(ctx).
		LoadData(client.StateKeyLaunches)
}
func (w *WorldLaunchPage) Render() app.UI {
	return &MainLayout{
		Content: []app.UI{
			components.AppLaunchItemList(),
		},
	}
}
