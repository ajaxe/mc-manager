package components

import (
	"github.com/ajaxe/mc-manager/internal/client"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type confirmModalData struct {
	title           string
	message         string
	show            bool
	confirmCallback app.EventHandler
}

type ConfirmModal struct {
	app.Compo
	title   string
	message string
	show    bool
	confirm app.EventHandler
}

func (c *ConfirmModal) OnMount(ctx app.Context) {
	ctx.Handle(client.ActionShowConfirm, c.handleConfirm)
}
func (c *ConfirmModal) Render() app.UI {
	return app.If(c.show, func() app.UI { return c.modal() }).
		Else(func() app.UI { return app.Div() })
}
func (c *ConfirmModal) modal() app.UI {
	return app.Div().
		Body(
			app.Div().Class("modal-backdrop fade show").Style("display", "block"),
			app.Div().Class("modal fade show").
				Style("display", "block").
				TabIndex(-1).Body(
				app.Div().
					Class("modal-dialog modal-dialog-centered").
					Body(
						app.Div().
							Class("modal-content").
							Body(
								app.Div().
									Class("modal-header").
									Body(
										app.H5().Class("modal-title").Text(c.title),
										app.Button().Type("button").Class("btn-close").
											DataSet("bs-dismiss", "modal").
											Aria("label", "Close").
											OnClick(c.close),
									),
								app.Div().Class("modal-body").
									Body(
										app.P().Text(c.message),
									),
								app.Div().
									Class("modal-footer").
									Body(
										app.Button().Type("button").Class("btn btn-secondary").
											DataSet("bs-dismiss", "modal").
											Aria("label", "Close").
											Text("No").
											OnClick(c.close),
										app.Button().Type("button").Class("btn btn-primary").
											Text("Yes").
											OnClick(c.confirmEventHandler),
									),
							),
					),
			),
		)
}
func (c *ConfirmModal) handleConfirm(ctx app.Context, a app.Action) {
	val, ok := a.Value.(confirmModalData)
	app.Logf("confirm action. ok: %v data: %v", ok, val)
	if ok {
		c.message = val.message
		c.title = val.title
		c.show = val.show
		c.confirm = val.confirmCallback
	}
}
func (c *ConfirmModal) close(ctx app.Context, e app.Event) {
	c.show = false
	c.title = ""
	c.message = ""
}
func (c *ConfirmModal) confirmEventHandler(ctx app.Context, e app.Event) {
	if c.confirm != nil {
		c.confirm(ctx, e)
	}
	c.close(ctx, e)
}
