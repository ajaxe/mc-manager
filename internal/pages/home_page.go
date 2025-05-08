package pages

import (
	"github.com/ajaxe/mc-manager/internal/components"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type HomePage struct {
	app.Compo
}

func (h *HomePage) Render() app.UI {
	return app.Main().
		Body(
			app.Div().Text("Welcome to the Home Page!").Class("home-page"),
			components.AppCodeUpdate(),
		)
}
