package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type PlayTimerActive struct {
	app.Compo
}

func (pt *PlayTimerActive) Render() app.UI {
	return app.Div().Class("card mt-2 bg-dark-subtle").Body(
		app.Div().Class("card-body row").Body(
			app.Div().Class("p-2 col-12").Body(
				app.Div().Class("h5 card-title").Text("Active Play Timer"),
			),
			app.Div().Class("p-2 col-sm-12 col-md-6").Body(
				app.Div().Class("form-text").ID("inp-min-help").
					Text("Enter the number of minutes to set the timer."),
			),
			app.Div().Class("p-2 col-sm-12 col-md-6").Body(
				app.Button().Class("btn btn-primary").Text("Start Timer"),
				app.Button().Class("btn btn-secondary ms-2").Text("Clear Timer"),
			),
		),
	)
}
