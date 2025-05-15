package components

import (
	"fmt"

	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type WorldItemCard struct {
	app.Compo
	item         *models.WorldItem
	active       bool
	intiGamemode string
}

func (w *WorldItemCard) Render() app.UI {
	return app.Div().
		ID(w.item.ID.Hex()).
		Class("card mt-2").Body(
		app.Div().Class("card-body").Body(
			app.H5().Class("card-title").
				Body(
					app.Text(w.item.Name+"  "),
					&WorldSelectBtn{
						active: w.active,
					},
				),
			app.H6().Class("card-subtitle mb-2 text-body-secondary").Text(w.item.Description),
			app.P().Class("card-text").Text("World Seed: "+w.item.WorldSeed),
			app.P().Class("card-text").
				Body(
					w.modeSelector(),
				),
		),
	)
}
func (w *WorldItemCard) modeSelector() app.UI {
	w.intiGamemode = w.item.GameMode
	return app.Span().Body(
		&FormSelect{
			Label: "World gamemode",
			SelectItems: map[string]string{
				"survival":  "Survival",
				"creative":  "Creative",
				"adventure": "Adventure",
			},
			Value:  w.item.GameMode,
			BindTo: &w.item.GameMode,
		},
		
		app.Button().Class("btn btn-link").Text("Change mode").
			OnClick(func(ctx app.Context, e app.Event) {
				e.PreventDefault()
				fmt.Println("Change mode: ", w.item.GameMode)
			}),
	)
}
