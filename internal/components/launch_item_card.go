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
					app.Div().Class("row").Body(
						app.Div().
							Class("col").
							Text(l.Item.Name),
						app.Div().
							Class("col").
							Text(l.Item.GameMode),
						app.Div().
							Class("col-3").
							Text(l.Item.WorldSeed),
						app.Div().
							Class("col").
							Text(l.Item.LaunchDate),
						app.Div().
							Class("col").
							Text(l.Item.Status),
					),
				),
		)
}
