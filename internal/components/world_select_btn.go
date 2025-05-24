package components

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type WorldSelectBtn struct {
	app.Compo
	active bool
}

func (w *WorldSelectBtn) Render() app.UI {
	c := "text-success"
	ico := "bi-check2-circle"
	if !w.active {
		c = "text-secondary"
		ico = "bi-circle"
	}
	return app.Span().
		Class("float-end badge rounded-pill").
		Body(
			app.I().
				Style("font-size", "1.8rem").
				Class(fmt.Sprintf("bi %s %s", ico, c)),
		)
}
