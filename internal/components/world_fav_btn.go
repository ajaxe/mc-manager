package components

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type WorldFavBtn struct {
	app.Compo
	color string
	ico   string
}

func (w *WorldFavBtn) Render() app.UI {
	if w.color == "" {
		w.color = "text-secondary"
		w.ico = "bi-star"
	}

	return app.Span().
		Class("float-start rounded-pill me-2").
		Body(
			app.Button().
				Class("btn btn-dark p-0 px-1").
				Type("button").
				Body(
					app.I().
						Style("font-size", "1.8rem").
						Class(fmt.Sprintf("bi %s %s", w.color, w.ico)),
				).
				OnClick(func(ctx app.Context, e app.Event) {
					w.toggleState()
				}),
		)
}
func (w *WorldFavBtn) toggleState() {
	if w.color == "text-secondary" {
		w.color = "text-warning"
		w.ico = "bi-star-fill"
	} else {
		w.color = "text-secondary"
		w.ico = "bi-star"
	}
}
