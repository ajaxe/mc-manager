package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type PlayTimerAdd struct {
	app.Compo
}

func (pt *PlayTimerAdd) Render() app.UI {
	return app.Div().Class("card mt-2 bg-dark-subtle").Body(
		app.Div().Class("card-body row").Body(
			app.Div().Class("p-2 col-12").Body(
				app.Div().Class("h5 card-title").Text("Play Timer"),
			),
			app.Div().Class("p-2 col-sm-12 col-md-6").Body(
				app.Input().
					Type("number").
					ID("inputNumber5").
					Class("form-control").
					Placeholder("Enter number of minutes").
					Aria("describedby", "inp-min-help"),
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
