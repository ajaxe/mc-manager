package components

import (
	"github.com/ajaxe/mc-manager/internal/client"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type PlayTimerActive struct {
	app.Compo
	ActiveTimer *models.PlayTimerItem
}

func (pt *PlayTimerActive) Render() app.UI {
	return app.Div().Class("card mt-2 bg-dark-subtle").Body(
		app.Div().Class("card-body row").Body(
			app.Div().Class("p-2 col-12").Body(
				app.Div().Class("h5 card-title").Text("Active Play Timer"),
			),
			app.Div().Class("p-2 col-sm-12 col-md-6").Body(
				app.Div().Body(
					app.Div().Class("form-text").
						Text("Play time ends at:"),
					app.Div().Class("form-text").
						Text(pt.formatDate(pt.ActiveTimer.EndDate)),
				),
				app.Div().Class("mt-2").Body(
					app.Button().Class("btn btn-primary").Text("Stop Timer").
						OnClick(func(ctx app.Context, e app.Event) {
							ctx.Async(func() {
								e := client.StopPlaytimer()
								ctx.Dispatch(func(ctx app.Context) {
									nctx := client.NewAppContext(ctx)
									nctx.ShowErrorMessage(nil, e)
									nctx.LoadData(client.StateKeyCurrentPlayTimer)
								})
							})
						}),
				),
			),
			app.Div().Class("p-2 col-sm-12 col-md-6").Body(
				&PlayTimerCountDown{EndDate: pt.ActiveTimer.EndDate},
			),
		),
	)
}
func (pt *PlayTimerActive) formatDate(s string) string {
	return client.BrowserDateDisplay(s)
}
