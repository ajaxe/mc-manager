package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type PlayTimer struct {
	app.Compo
}

func (pt *PlayTimer) Render() app.UI {
	return &PlayTimerAdd{}
}
