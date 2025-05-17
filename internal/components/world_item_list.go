package components

import (
	"fmt"

	"github.com/ajaxe/mc-manager/internal/client"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type WorldItemList struct {
	app.Compo
	items []*models.WorldItem
}

func (w *WorldItemList) OnMount(ctx app.Context) {
	fmt.Printf("component mounted: %s\n", app.Window().URL())
	ctx.ObserveState(client.StateKeyWorlds, &w.items).
		OnChange(func() {
			app.Log("OnMount world list")
		})
}

func (w *WorldItemList) Render() app.UI {

	return app.Div().
		Class("row row-cols-1 row-cols-md-2 g-4").
		Body(
			app.Range(w.items).
				Slice(func(i int) app.UI {
					return app.Div().Class("col").Body(
						&WorldItemCard{
							item:   w.items[i],
							active: i == 0,
						},
					)
				}),
		)
}
