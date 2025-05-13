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
	e := []app.UI{ m.appName(), components.AppCodeUpdate()}

	e = append(e, m.Content...)

	return app.Div().Class("container").Body(e...)
}

func (m *MainLayout) appName() app.UI {
	return app.Div().Class("row").Body(
		app.Div().Class("col").Body(
			components.NewAppName(),
		),
	)
}
