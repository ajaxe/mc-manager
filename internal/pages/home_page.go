package pages

import (
	"github.com/ajaxe/mc-manager/internal/client"
	"github.com/ajaxe/mc-manager/internal/components"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type HomePage struct {
	app.Compo
}

func (h *HomePage) OnNav(ctx app.Context) {
	client.NewAppContext(ctx).
		LoadData(client.StateKeyWorlds)
}
func (h *HomePage) Render() app.UI {
	return &MainLayout{
		Content: []app.UI{
			components.AppWorldItemList(),
		},
	}
}
