package components

import (
	"strconv"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type PlayTimerCountDown struct {
	app.Compo
	Minutes int
}

func (pt *PlayTimerCountDown) Render() app.UI {
	return app.Div().
		Body(
			app.Div().Text("Time remaining:"),
			app.Div().Body(
				app.Span().Class("display-1").Text(strconv.Itoa(pt.Minutes)),
				app.Span().Text("min"),
				app.Span().Class("ms-2 display-1").Text("00"),
				app.Span().Text("sec"),
			),
		)
}
