package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type AppName struct {
	app.Compo
}

func (a *AppName) Render() app.UI {
	return app.Div().Class("d-flex").Body(
		a.logo(),
		app.Div().Class("app-name").Body(
			app.H1().Class("d-none d-md-block").Text("Minecraft Server Manager"),
			app.H1().Class("d-sm-none").Text("MC Manager"),
			app.P().Text("Manage your Minecraft servers with ease."),
		),
	)
}

func (a *AppName) logo() app.UI {
	return app.Div().
		Class("p-2").
		Body(
			app.Img().
				Src("/web/favicon.png").
				Class("app-icon").
				Alt("Minecraft Server Manager Icon").
				Width(64).
				Height(64),
		)
}
