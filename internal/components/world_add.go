package components

import (
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type WorldAdd struct {
	app.Compo
	WordItem *models.WorldItem
	addMode  bool
}

func (w *WorldAdd) Render() app.UI {
	return app.Div().Class("row").Body(
		app.Div().Class("col").Body(
			app.If(!w.addMode, func() app.UI {
				return w.addBtn()
			}),
			app.If(w.addMode, func() app.UI {
				return w.form()
			}),
		),
	)
}

func (w *WorldAdd) addBtn() app.UI {
	c := "btn-primary"
	t := "Add a new world"
	return app.Div().Body(
		app.Button().
			Class("btn " + c).
			Text(t).
			OnClick(w.showAddView),
	)
}

func (w *WorldAdd) form() app.UI {
	wi := &models.WorldItem{}
	return app.Div().Class("card mt-2").Body(
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
				Text("Cancel").
				OnClick(func(ctx app.Context, e app.Event) {
					w.addMode = false
				}),
		),
	)
}

func (w *WorldAdd) showAddView(ctx app.Context, e app.Event) {
	if w.addMode {
		w.addMode = false
	} else {
		w.addMode = true
	}
}
