package components

import (
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
					app.H3().
						Class("card-title").
						Text(l.Item.Name),
					app.P().
						Class("card-text").
						Text(l.Item.LaunchDate),
					app.Button().
						Class("btn btn-primary").
						Text("Launch"),
				),
		)
}
