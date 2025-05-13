package components

import (
	"fmt"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type FormLabel struct {
	app.Compo
	For   string
	Label string
}

func (l *FormLabel) Render() app.UI {
	return app.Label().
		Class("form-label").
		For(l.For).
		Text(l.Label)
}

type FormText struct {
	app.Compo
	ID          string
	Value       string
	ReadOnly    bool
	BindTo      any
	InputType   string
	Placeholder string
}

func (t *FormText) Render() app.UI {
	c := "form-control"
	if t.ReadOnly {
		c += "-plaintext"
	}
	it := t.InputType
	if it == "" {
		it = "text"
	}
	if t.ID == "" {
		t.ID = fmt.Sprintf("txt-%v", time.Now().UnixMicro())
	}

	elem := app.Input().
		Type(it).
		ReadOnly(t.ReadOnly).
		Placeholder(t.Placeholder).
		Class(c).
		ID(t.ID).
		Value(t.Value)

	if !t.ReadOnly {
		elem.OnChange(t.ValueTo(t.BindTo))
	}

	return elem
}

type FormControl struct {
	app.Compo
	Content []app.UI
	Compact bool
}

func (f *FormControl) Render() app.UI {
	m := "mb-3"
	if f.Compact {
		m = ""
	}
	return app.Div().Class("form-floating " + m).Body(f.Content...)
}

type FormCheckbox struct {
	app.Compo
	role     string // empty or "switch"
	id       string
	BindTo   *bool
	Value    bool
	label    string
	Disabled bool
}

func (c *FormCheckbox) Render() app.UI {
	s := ""
	if c.role == "switch" {
		s = "form-switch"
	}

	if c.id == "" {
		c.id = fmt.Sprintf("chk-%v", time.Now().UnixMicro())
	}

	input := app.Input().
		Class("form-check-input").
		Type("checkbox").
		Checked(c.Value).
		Disabled(c.Disabled).
		ID(c.id)

	if c.BindTo != nil {
		input.OnChange(c.toggleCheck(c.BindTo))
	}

	return app.Div().Class("form-check "+s).
		Body(
			input,
			app.Label().Class("form-check-label").
				For(c.id).
				Text(c.label),
		)
}

func (c *FormCheckbox) toggleCheck(v *bool) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		checked := ctx.JSSrc().Get("checked").Bool()
		*v = checked
	}
}

type FormSelect struct {
	app.Compo
	Label       string
	SelectItems map[string]string
	Value       string
	BindTo      any
}

func (f *FormSelect) Render() app.UI {
	return app.Select().
		Class("form-select").
		Aria("label", f.Label).
		Body(
			app.Option().Selected(true).Text(f.Label),
			app.Range(f.SelectItems).Map(func(key string) app.UI {
				return app.Option().
					Value(key).
					Text(f.SelectItems[key]).
					Selected(f.Value == key)
			}),
		).
		OnChange(f.ValueTo(f.BindTo))
}
