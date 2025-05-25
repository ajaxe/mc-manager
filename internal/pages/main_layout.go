package pages

import (
	"github.com/ajaxe/mc-manager/internal/client"
	"github.com/ajaxe/mc-manager/internal/components"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type MainLayout struct {
	app.Compo
	Content []app.UI
}

func (m *MainLayout) OnNav(ctx app.Context) {
	ctx.Async(func() {
		client.LoginCheck()
	})
}

func (m *MainLayout) Render() app.UI {

	return app.Div().Class("container").Body(
		m.appName(),
		components.AppCodeUpdate(),
		components.AppConfirmModal(),
		app.Div().Class("row mt-4").Body(
			app.Div().Class("col-auto").
				Body(
					components.AppSidebarMenu(),
				),
			app.Div().Class("col").Body(
				app.Main().Body(m.Content...),
			),
		),
	)
}

func (m *MainLayout) appName() app.UI {
	return app.Div().Class("row").Body(
		app.Div().Class("col").Body(
			components.NewAppName(),
		),
	)
}
