package pages

import (
	"github.com/ajaxe/mc-manager/internal/components"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type WorldAddPage struct {
	app.Compo
}

func (w *WorldAddPage) Render() app.UI {
	return &MainLayout{
		Content: []app.UI{
			components.AppWorldAdd(),
		},
	}
}
