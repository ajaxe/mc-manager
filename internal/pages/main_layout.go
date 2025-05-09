package pages

import (
	"github.com/ajaxe/mc-manager/internal/components"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type MainLayout struct {
	app.Compo
	Content []app.UI
}

func (m *MainLayout) Render() app.UI {
	return app.Div().Class("container").Body(
		app.Div().Class("row").Body(
			app.Div().Class("col").Body(
				components.NewAppName(),
			),
		),
		app.Div().Class("row").Body(
			app.Div().Class("col").Body(
				m.Content...,
			),
		),
		components.AppCodeUpdate(),
	)
}
