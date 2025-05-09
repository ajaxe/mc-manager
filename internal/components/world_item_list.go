package components

import (
	"fmt"

	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type WorldItemList struct {
	app.Compo
	items []*models.WorldItem
}

func (w *WorldItemList) OnNav(ctx app.Context) {
	fmt.Printf("component navigated: %s\n", app.Window().URL())
	w.items = []*models.WorldItem{
		{
			ID:          bson.NewObjectID(),
			Name:        "World 1",
			Description: "Description for World 1",
			WorldSeed:   "Seed 1",
		},
		{
			ID:          bson.NewObjectID(),
			Name:        "World 2",
			Description: "Description for World 2",
			WorldSeed:   "Seed 2",
		},
	}
}
func (w *WorldItemList) OnMount(ctx app.Context) {

}

func (w *WorldItemList) Render() app.UI {
	return app.Div().Class("row").Body(
		app.Div().Class("col").Body(
			app.Range(w.items).Slice(func(i int) app.UI {
				return &WorldItemCard{
					Item: w.items[i],
				}
			}),
		),
	)
}
