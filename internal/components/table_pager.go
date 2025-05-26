package components

import (
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type TablePager struct {
	app.Compo
	PrevID  string
	NextID  string
	OnClick func(app.Context, string)
}

func (t *TablePager) Render() app.UI {
	prev := ""
	if t.PrevID == "" {
		prev = "disabled"
	}
	next := ""
	if t.NextID == "" {
		next = "disabled"
	}
	return app.Nav().
		Body(
			app.Ul().Class("pagination justify-content-end").
				Body(
					t.linkItem("Previous", prev, models.PageDirectionPrev),
					t.linkItem("Next", next, models.PageDirectionNext),
				),
		)
}

func (t *TablePager) linkItem(text, disabled, val string) app.UI {
	a := app.A().Class("page-link").Text(text).Href("#").OnClick(func(ctx app.Context, e app.Event) {
		t.OnClick(ctx, val)
	})
	return app.Li().Class("page-item " + disabled).
		Body(a)
}
