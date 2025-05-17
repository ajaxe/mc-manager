package components

import (
	"github.com/ajaxe/mc-manager/internal/client"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type WorldAdd struct {
	app.Compo
	WordItem  *models.WorldItem
	message   string
	isSuccess bool
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
				Text("Add").
				OnClick(func(ctx app.Context, e app.Event) {
					e.PreventDefault()
					ctx.Async(func() {
						r, e := client.WorldCreate(wi)
						ctx.Dispatch(func(ctx app.Context) {
							if e != nil {
								w.message = e.Error()
								w.isSuccess = false
							} else if !r.Success {
								w.message = r.ErrorMessage
								w.isSuccess = false
							} else {
								w.message = "World added successfully"
								w.isSuccess = true
							}
							ctx.Update()
						})
					})
				}),
			app.Button().
				Class("btn btn-secondary ms-2").
				Text("Cancel").
				OnClick(func(ctx app.Context, e app.Event) {
					e.PreventDefault()
					ctx.Navigate("/")
				}),
			app.If(w.hasMessage(!w.isSuccess), func() app.UI {
				return app.Span().Class("text-danger ms-3 fw-semibold").Text(w.message)
			}),
			app.If(w.isSuccess, func() app.UI {
				return app.Span().Class("text-success ms-3 fw-semibold").Text(w.message)
			}),
		),
	)
}

func (w *WorldAdd) hasMessage(valid bool) bool {
	return w.message != "" && valid
}
