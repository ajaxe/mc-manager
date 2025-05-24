package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type WorldDeleteBtn struct {
	app.Compo
	active bool
}

func (w *WorldDeleteBtn) Render() app.UI {
	c := "text-danger"
	if w.active {
		c = "text-secondary"
	}
	return app.Span().
		Class("float-end rounded-pill").
		Body(
			app.If(w.active, func() app.UI {
				return app.I().
					Style("font-size", "1.8rem").
					Class("bi bi-x-lg " + c)
			}).Else(func() app.UI {
				return app.Button().
					Class("btn btn-dark p-0 px-1").
					Type("button").
					Body(
						app.I().
							Style("font-size", "1.8rem").
							Class("bi bi-x-lg " + c),
					)
			}),
		)
}
