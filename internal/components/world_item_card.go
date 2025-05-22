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
	disabled     bool
	loadMessage  string
}

func (w *WorldItemCard) Render() app.UI {
	b := ""
	if w.Item.IsActive {
		b = "border-success"
	}
	return app.Div().
		ID(w.Item.ID).
		Class("card mt-2 bg-dark-subtle "+b).
		Body(
			&CardSpinner{
				Show:    w.disabled,
				Message: w.loadMessage,
			},
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
					app.Label().For("ws-"+w.Item.ID).Text("World Seed: "),
					app.Span().ID("ws-"+w.Item.ID).Class("ms-3").Text(w.Item.WorldSeed),
				),
				app.P().Class("card-text").
					Body(
						w.modeSelector(),
					),

				w.launchWorldBtn(),

				app.If(!w.Item.IsActive, func() app.UI {
					return app.Button().Class("btn btn-link").Disabled(w.disabled).Text("Delete world").
						OnClick(func(ctx app.Context, e app.Event) {
							e.PreventDefault()
							ctx.Async(func() {
								client.WorldDelete(w.Item.ID)
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
	id := fmt.Sprintf("select-%s", w.Item.ID)
	return app.Span().Body(
		&FormLabel{
			For:   id,
			Label: "World gamemode",
		},
		&FormSelect{
			ID:       id,
			Label:    "World gamemode",
			Disabled: w.disabled,
			SelectItems: map[string]string{
				"survival":  "Survival",
				"creative":  "Creative",
				"adventure": "Adventure",
			},
			Value:  w.Item.GameMode,
			BindTo: &w.Item.GameMode,
			OnChange: func(ctx app.Context, e app.Event) {
				w.Item.GameMode = ctx.JSSrc().Get("value").String()
				w.disabled = true
				w.loadMessage = "Updating game mode."
				//ctx.NewActionWithValue(client.ActionShowCardSpinners, true)
				ctx.Async(func() {
					_, _ = client.WorldUpdate(w.Item)

					ctx.Dispatch(func(ctx app.Context) {
						w.disabled = false
						w.loadMessage = ""
					})
				})
			},
		},
	)
}
func (w *WorldItemCard) launchWorldBtn() app.UI {
	txt := "Launch world"
	m := "Launching world."
	if w.Item.IsActive {
		txt = "Change game mode"
		m = "Re-launching world in new Game-mode"
	}
	return app.Button().Class("btn btn-link").Disabled(w.disabled).Text(txt).
		OnClick(func(ctx app.Context, e app.Event) {
			e.PreventDefault()
			w.disabled = true
			w.loadMessage = m
			ctx.NewActionWithValue(client.ActionShowCardSpinners, true)
			ctx.Async(func() {
				_ = client.LaunchWorld(w.Item)
				// TODO: erorr handling
				ctx.Dispatch(func(ctx app.Context) {
					ctx.NewActionWithValue(client.ActionShowCardSpinners, false)
					w.disabled = false
					w.loadMessage = ""
					client.NewAppContext(ctx).
						LoadData(client.StateKeyWorlds)
				})
			})
		})
}
