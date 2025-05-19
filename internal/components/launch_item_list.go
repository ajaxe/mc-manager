package components

import (
	"fmt"

	"github.com/ajaxe/mc-manager/internal/client"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type LaunchItemList struct {
	app.Compo
	items []*models.LaunchItem
}

func (l *LaunchItemList) OnMount(ctx app.Context) {
	fmt.Printf("component mounted: %s\n", app.Window().URL())
	ctx.ObserveState(client.StateKeyLaunches, &l.items)
}
func (l *LaunchItemList) Render() app.UI {
	return app.Div().
		//Class("row g-4").
		Body(
			l.header(),
			app.Range(l.items).
				Slice(func(i int) app.UI {
					return &LaunchItemCard{
						Item: l.items[i],
					}
				}),
		)
}

func (l *LaunchItemList) header() app.UI {
	return app.Div().
		Class("card d-none d-md-block").
		Body(
			app.Div().
				Class("card-body fw-bold").
				Body(
					app.Div().Class("row").Body(
						app.Div().
							Class("col").
							Text("World Name"),
						app.Div().
							Class("col text-capitalize").
							Text("GameMode"),
						app.Div().
							Class("col-3").
							Text("World Seed"),
						app.Div().
							Class("col").
							Text("Launch Date"),
						app.Div().
							Class("col").
							Text("Status"),
					),
				),
		)
}
