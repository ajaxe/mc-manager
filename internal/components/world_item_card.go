package components

import (
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type WorldItemCard struct {
	app.Compo
	Item *models.WorldItem
}

func (w *WorldItemCard) Render() app.UI {
	return app.Div().
		ID(w.Item.ID.Hex()).
		Class("card mt-2").Body(
		app.Div().Class("card-body").Body(
			app.H5().Class("card-title").Text(w.Item.Name),
			app.H6().Class("card-subtitle mb-2 text-body-secondary").Text(w.Item.Description),
			app.P().Class("card-text").Text("WorldSeed: "+w.Item.WorldSeed),
		),
	)
}
