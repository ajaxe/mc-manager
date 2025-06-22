package components

import (
	"fmt"

	"github.com/ajaxe/mc-manager/internal/client"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type LaunchItemCard struct {
	app.Compo
	Item *models.LaunchItem
}

func (l *LaunchItemCard) Render() app.UI {
	return app.Div().
		Class("card").
		Body(
			app.Div().
				Class("card-body").
				Body(
					app.Div().Class("row").Body(
						l.displayField(l.asText(l.Item.Name), "col-sm-12 col-md", "World Name:"),
						l.displayField(l.asText(l.Item.GameMode), "col-sm-12 col-md text-capitalize", "Gamemode:"),
						l.displayField(l.asText(l.Item.WorldSeed), "col-sm-12 col-md-3", "World Seed:"),
						l.displayField(l.asText(l.dtDisplay()), "col-sm-12 col-md", "Launch Date:"),
						l.displayField(l.statusIcon(), "col-sm-12 col-md", "Status:"),
					),
				),
		)
}
func (l *LaunchItemCard) asText(s string) app.UI {
	return app.Text(s)
}
func (l *LaunchItemCard) displayField(el app.UI, css, lbl string) app.UI {
	return app.Div().
		Class(css+" pb-2 pb-md-0").
		Body(
			app.Label().
				Class("d-sm-inline-block d-md-none pe-3 fw-medium").
				Style("min-width", "110px").
				Text(lbl),
			el,
		)

}
func (l *LaunchItemCard) statusIcon() app.UI {
	ico := "bi-check-circle-fill"
	co := "text-success"

	if l.Item.Status != "success" {
		ico = "bi-exclamation-octagon-fill"
		co = "text-danger"
	}
	return app.Span().
		Body(
			app.I().
				Class(fmt.Sprintf("bi %s %s", ico, co)),
			app.Span().
				Class("ms-2 text-capitalize").
				Text(l.Item.Status),
		)
}
func (l *LaunchItemCard) dtDisplay() string {
	return client.BrowserDateDisplay(l.Item.LaunchDate)
}
