package components

import (
	"github.com/ajaxe/mc-manager/internal/client"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type CardSpinner struct {
	app.Compo
	Show       bool
	showGlobal bool
	Message    string
}

func (c *CardSpinner) OnMount(ctx app.Context) {
	ctx.Handle(client.ActionShowCardSpinners, func(ctx app.Context, a app.Action) {
		c.showGlobal = a.Value.(bool)
	})
}

func (c *CardSpinner) Render() app.UI {
	css := "d-none"
	if c.Show || c.showGlobal {
		css = ""
	}
	return app.Div().Class("static-blocker rounded d-flex " + css).
		Body(
			app.Div().Class("align-self-end ms-auto p-2").
				Body(
					app.Div().Class("float-end").Body(
						app.Div().Class("spinner-border").Role("status").
							Body(
								app.Span().Class("visually-hidden").Text("Loading..."),
							),
					),
					app.If(c.Message != "", func() app.UI {
						return app.Div().Style("clear", "both").
							Text(c.Message)
					}),
				),
		)
}
