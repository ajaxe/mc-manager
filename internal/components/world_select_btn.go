package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type WorldSelectBtn struct {
	app.Compo
	active bool
}

func (w *WorldSelectBtn) Render() app.UI {
	c := "text-success"
	if !w.active {
		c = "text-secondary"
	}
	return app.Span().
		Class("float-end badge rounded-pill").
		Body(
			app.I().
				Style("font-size", "1.8rem").
				Class("ms-2 bi bi-check2-square " + c),
		)
}
