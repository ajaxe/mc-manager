package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type CodeUpdate struct {
	app.Compo

	updateAvailable bool
}

func (c *CodeUpdate) OnAppUpdate(ctx app.Context) {
	c.updateAvailable = ctx.AppUpdateAvailable()
}

func (c *CodeUpdate) Render() app.UI {
	return app.Div().Class("toast-container top-0 start-50 translate-middle-x").
		Body(
			app.If(c.updateAvailable, func() app.UI {
				return c.updateToast()
			}),
		)
}

func (c *CodeUpdate) updateToast() app.UI {
	return app.Div().
		ID("app-update-toast").
		Class("toast align-items-center show").
		Role("alert").
		Aria("live", "assertive").
		Aria("atomic", "true").
		Body(
			app.Div().Class(("d-flex")).Body(
				app.Div().Class("toast-body").Body(
					app.P().Class("fw-normal").Body(
						app.Text("A new version of the app is available."),
					),
				),
				app.Button().
					Class("btn btn-primary btn-sm me-2 m-auto").
					Text("Update").
					OnClick(c.onUpdateClick),
				app.Button().
					Class("btn-close me-2 m-auto").
					Aria("label", "Close").
					DataSet("bs-dismiss", "toast"),
			),
		)
}

func (a *CodeUpdate) onUpdateClick(ctx app.Context, e app.Event) {
	// Reloads the page to display the modifications.
	ctx.Reload()
}
