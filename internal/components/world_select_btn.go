package components

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type WorldSelectBtn struct {
	app.Compo
	active   bool
	OnSelect app.EventHandler
}

func (w *WorldSelectBtn) Render() app.UI {
	c := "text-success"
	ico := "bi-check2-circle"
	if !w.active {
		c = "text-secondary"
		ico = "bi-circle"
	}
	return app.Span().
		Class("float-end rounded-pill ms-1").
		Body(
			app.If(w.active, func() app.UI {
				return app.I().
					Style("font-size", "1.8rem").
					Class(fmt.Sprintf("bi %s %s", ico, c))
			}).Else(func() app.UI {
				return app.Button().
					Class("btn btn-dark p-0 px-1").
					Type("button").
					Body(
						app.I().
							Style("font-size", "1.8rem").
							Class(fmt.Sprintf("bi %s %s", ico, c)),
					).
					OnClick(w.onClick)
			}),
		)
}
func (w *WorldSelectBtn) onClick(ctx app.Context, e app.Event) {
	if w.OnSelect != nil && !w.active {
		w.OnSelect(ctx, e)
	}
}
