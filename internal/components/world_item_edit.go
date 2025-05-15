package components

import (
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type WorldItemEdit struct {
	app.Compo
	WorldItem *models.WorldItem
	title     string
}

func (w *WorldItemEdit) Render() app.UI {
	return app.Div().
		Class("mb-2").
		Body(
			app.H5().Text(w.title),
			app.Form().Body(
				&FormControl{
					Compact: false,
					Content: []app.UI{
						&FormText{
							ID:          "world-name",
							Placeholder: "World name",
							Value:       w.WorldItem.Name,
						},
						&FormLabel{
							For:   "world-name",
							Label: "World name",
						},
					},
				},
				&FormControl{
					Compact: false,
					Content: []app.UI{
						&FormText{
							ID:          "world-desc",
							Placeholder: "World description",
							Value:       w.WorldItem.Description,
						},
						&FormLabel{
							For:   "world-desc",
							Label: "World description",
						},
					},
				},
				&FormControl{
					Compact: false,
					Content: []app.UI{
						&FormText{
							ID:          "world-seed",
							Placeholder: "World seed",
							Value:       w.WorldItem.WorldSeed,
						},
						&FormLabel{
							For:   "world-seed",
							Label: "World seed",
						},
					},
				},
				&FormSelect{
					Label: "World gamemode",
					SelectItems: map[string]string{
						"survival":  "Survival",
						"creative":  "Creative",
						"adventure": "Adventure",
					},
					Value:  w.WorldItem.GameMode,
					BindTo: &w.WorldItem.GameMode,
				},
			),
		)
}
