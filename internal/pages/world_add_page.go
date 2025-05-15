package pages

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type WorldAddPage struct {
	app.Compo
}

func (w *WorldAddPage) Render() app.UI {
	return &MainLayout{
		Content: []app.UI{
			app.H1().Text("World Add Page"),
			app.P().Text("This is the world add page."),
		},
	}
}
