package components

import (
	"strconv"

	"github.com/ajaxe/mc-manager/internal/client"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type PlayTimerAdd struct {
	app.Compo
	valueMinutes string
}

func (pt *PlayTimerAdd) Render() app.UI {
	input := app.Input().
		Type("number").
		ID("inputNumber5").
		Class("form-control").
		Placeholder("Enter number of minutes").
		Aria("describedby", "inp-min-help").
		Value(pt.valueMinutes).
		OnChange(pt.ValueTo(&pt.valueMinutes))

	return app.Div().Class("card mt-2 bg-dark-subtle").Body(
		app.Div().Class("card-body row").Body(
			app.Div().Class("p-2 col-12").Body(
				app.Div().Class("h5 card-title").Text("Play Timer"),
			),
			app.Div().Class("p-2 col-sm-12 col-md-6").Body(
				input,
				app.Div().Class("form-text").ID("inp-min-help").
					Text("Enter the number of minutes to set the timer."),
			),
			app.Div().Class("p-2 col-sm-12 col-md-6").Body(
				app.Button().Class("btn btn-primary").Text("Start Timer").
					Disabled(pt.startDisabled()).
					OnClick(func(ctx app.Context, e app.Event) {
						ctx.Async(func() {
							minutes, _ := strconv.Atoi(pt.valueMinutes)
							e := client.StartPlaytimer(&models.PlayTimerItem{
								Minutes: minutes,
							})
							if e != nil {
								client.NewAppContext(ctx).
									ShowErrorMessage(nil, e)
							} else {
								client.NewAppContext(ctx).LoadData(client.StateKeyCurrentPlayTimer)
							}
						})
					}),
				app.Button().Class("btn btn-secondary ms-2").Text("Clear Timer").
					OnClick(func(ctx app.Context, e app.Event) {
						pt.valueMinutes = "" // Clear the input field
						ctx.Update()
					}),
			),
		),
	)
}
func (pt *PlayTimerAdd) startDisabled() bool {
	min, err := strconv.Atoi(pt.valueMinutes)
	return min <= 0 || err != nil
}
