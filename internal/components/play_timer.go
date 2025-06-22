package components

import (
	"github.com/ajaxe/mc-manager/internal/client"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type PlayTimer struct {
	app.Compo
	activeTimer models.PlayTimerItem
}

func (pt *PlayTimer) OnMount(ctx app.Context) {
	ctx.ObserveState(client.StateKeyCurrentPlayTimer, &pt.activeTimer)
}

func (pt *PlayTimer) Render() app.UI {
	return app.If(pt.timerActive(), func() app.UI { return &PlayTimerActive{ActiveTimer: &pt.activeTimer} }).
		Else(func() app.UI { return &PlayTimerAdd{} })
}

func (pt *PlayTimer) timerActive() bool {
	return pt.activeTimer.IsActive
}
