package components

import (
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type WorldAdd struct {
	app.Compo
	WordItem *models.WorldItem
}

func (w *WorldAdd) Render() app.UI {
	return app.Div().Class("row").Body(
		app.Div().Class("col").Body(
			w.form(),
		),
	)
}

func (w *WorldAdd) form() app.UI {
	wi := &models.WorldItem{}
	return app.Div().Class("card mt-2 bg-dark-subtle ").Body(
		app.Div().Class("card-body").Body(
			&WorldItemEdit{
				WorldItem: wi,
				title:     "Add a new world",
			},
			app.Button().
				Class("btn btn-primary").
				Text("Add"),
			app.Button().
				Class("btn btn-secondary ms-2").
				Text("Cancel"),
		),
	)
}

