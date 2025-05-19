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
			app.Range(l.items).
				Slice(func(i int) app.UI {
					return &LaunchItemCard{
						Item: l.items[i],
					}
				}),
		)
}
