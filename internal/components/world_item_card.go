package components

import (
	"fmt"

	"github.com/ajaxe/mc-manager/internal/client"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type WorldItemCard struct {
	app.Compo
	Item         *models.WorldItem
	intiGamemode string
}

func (w *WorldItemCard) Render() app.UI {
	b := ""
	if w.Item.IsActive {
		b = "border-success"
	}
	return app.Div().
		ID(w.Item.ID.Hex()).
		Class("card mt-2 bg-dark-subtle " + b).Body(
		app.Div().Class("card-body").Body(
			app.H5().Class("card-title").
				Body(
					app.Text(w.Item.Name+"  "),
					&WorldSelectBtn{
						active: w.Item.IsActive,
					},
				),
			app.H6().Class("card-subtitle mb-2 text-body-secondary").Text(w.Item.Description),
			app.P().Class("card-text").Body(
				app.Label().For("ws-"+w.Item.ID.Hex()).Text("World Seed: "),
				app.Span().ID("ws-"+w.Item.ID.Hex()).Class("ms-3").Text(w.Item.WorldSeed),
			),
			app.P().Class("card-text").
				Body(
					w.modeSelector(),
				),
			app.If(w.Item.IsActive, func() app.UI {
				return app.Button().Class("btn btn-link").Text("Change mode").
					OnClick(func(ctx app.Context, e app.Event) {
						e.PreventDefault()
						fmt.Println("Change mode: ", w.Item.GameMode)
					})
			}),
			app.If(!w.Item.IsActive, func() app.UI {
				return app.Button().Class("btn btn-link").Text("Launch world").
					OnClick(func(ctx app.Context, e app.Event) {
						e.PreventDefault()
						ctx.Async(func() {
							_ = client.LaunchWorld(w.Item.ID)
							// TODO: erorr handling
							ctx.Dispatch(func(ctx app.Context) {
								client.NewAppContext(ctx).
									LoadData(client.StateKeyWorlds)
							})
						})
					})
			}),
			app.If(!w.Item.IsActive, func() app.UI {
				return app.Button().Class("btn btn-link").Text("Delete world").
					OnClick(func(ctx app.Context, e app.Event) {
						e.PreventDefault()
						ctx.Async(func() {
							client.WorldDelete(w.Item.ID.Hex())
							ctx.Dispatch(func(ctx app.Context) {
								client.NewAppContext(ctx).
									LoadData(client.StateKeyWorlds)
							})
						})
					})
			}),
		),
	)
}
func (w *WorldItemCard) modeSelector() app.UI {
	w.intiGamemode = w.Item.GameMode
	id := fmt.Sprintf("select-%s", w.Item.ID.Hex())
	return app.Span().Body(
		&FormLabel{
			For:   id,
			Label: "World gamemode",
		},
		&FormSelect{
			ID:    id,
			Label: "World gamemode",
			SelectItems: map[string]string{
				"survival":  "Survival",
				"creative":  "Creative",
				"adventure": "Adventure",
			},
			Value:  w.Item.GameMode,
			BindTo: &w.Item.GameMode,
		},
	)
}
