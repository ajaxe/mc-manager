package components

import (
	"github.com/ajaxe/mc-manager/internal/client"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type StatusToast struct {
	app.Compo
	status  string
	message string
}

func (s *StatusToast) OnMount(ctx app.Context) {
	ctx.Handle(client.ActionStatusToast, s.showToast)
}

func (s *StatusToast) Render() app.UI {
	return app.If(s.message == "", func() app.UI { return app.Div().Class("d-none") }).
		Else(func() app.UI { return s.toast() })
}

func (s *StatusToast) toast() app.UI {
	bg := "text-bg-secondary"
	if s.status == client.ToastStatusSuccess {
		bg = "text-bg-success"
	} else if s.status == client.ToastStatusError {
		bg = "text-bg-danger"
	}
	return app.Div().Class("toast-container p-3 position-fixed top-0 end-0").
		Body(
			app.Div().
				Class("toast align-items-center border-0 fade show "+bg).
				Role("alert").
				Aria("live", "assertive").
				Aria("atomic", "true").
				Body(
					app.Div().Class("d-flex").Body(
						app.Div().Class("toast-body").Body(
							app.Text(s.message),
						),
						app.Button().Type("button").
							Class("btn-close btn-close-white me-2 m-auto").
							DataSet("bs-dismiss", "toast").
							Aria("label", "close"),
					),
				),
		)
}

func (s *StatusToast) showToast(ctx app.Context, a app.Action) {
	d, ok := a.Value.(client.StatusToastData)
	if !ok {
		return
	}
	app.Logf("showToast: m:%s, s:%s", d.Message, d.Status)
	s.message = d.Message
	s.status = d.Status
	app.Window().Call("setTimeout", app.FuncOf(func(this app.Value, args []app.Value) any {
		s.message = ""
		ctx.Update()
		return nil
	}), 5000)
}
