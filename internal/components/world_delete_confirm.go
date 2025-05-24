package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type WorldDeleteConfirm struct {
	app.Compo
	Show bool
}

func (w *WorldDeleteConfirm) Render() app.UI {
	css := "d-none"
	if w.Show {
		css = ""
	}
	return app.Div().Class("static-blocker rounded d-flex " + css).
		Body(
			app.Div().Class("align-self-end ms-auto p-2").
				Body(
					app.Div().Text("confirm comes here"),
				),
		)
}
