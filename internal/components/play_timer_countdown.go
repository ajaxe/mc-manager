package components

import (
	"strconv"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type PlayTimerCountDown struct {
	app.Compo
	EndDate string
	min     int
	sec     int
	dt      *time.Time
	timer   app.Value
	ctx     app.Context
}

func (pt *PlayTimerCountDown) OnDismount() {
	if pt.timer != nil {
		app.Window().Call("clearInterval", pt.timer)
		pt.timer = nil
	}
	pt.dt = nil
	pt.min = 0
	pt.sec = 0
}
func (pt *PlayTimerCountDown) OnMount(ctx app.Context) {
	pt.ctx = ctx
}

func (pt *PlayTimerCountDown) Render() app.UI {
	if pt.dt == nil {
		t, _ := time.Parse(time.RFC3339, pt.EndDate)
		pt.dt = &t
		pt.initializeTimer()
	}
	return app.Div().
		Body(
			app.Div().Text("Time remaining:"),
			app.Div().Body(
				app.Span().Class("display-1").Text(strconv.Itoa(pt.min)),
				app.Span().Text("min"),
				app.Span().Class("ms-2 display-1").Text(strconv.Itoa(pt.sec)),
				app.Span().Text("sec"),
			),
		)
}
func (pt *PlayTimerCountDown) initializeTimer() {
	pt.timer = app.Window().Call("setInterval", app.FuncOf(func(this app.Value, args []app.Value) any {
		now := time.Now()
		diff := pt.dt.Sub(now)
		if diff <= 0 {
			pt.min = 0
			pt.sec = 0
			app.Window().Call("clearInterval", pt.timer)
			app.Logf("Play timer has ended!")
			return nil
		}
		pt.min = int(diff.Minutes())
		pt.sec = int(diff.Seconds()) % 60
		pt.ctx.Update()
		return nil
	}), 1000)
}
