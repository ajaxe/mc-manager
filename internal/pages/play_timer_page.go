package pages

import (
	"github.com/ajaxe/mc-manager/internal/components"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type PlayTimerPage struct {
	app.Compo
}

func (p *PlayTimerPage) Render() app.UI {
	return &MainLayout{
		Content: []app.UI{
			components.AppPlayTimer(),
		},
	}
}
