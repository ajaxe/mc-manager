package pages

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type WorldLaunchPage struct {
	app.Compo
}

func (w *WorldLaunchPage) Render() app.UI {
	return &MainLayout{
		Content: []app.UI{
			app.H1().Text("World Launch Page"),
			app.P().Text("This is the world launch page."),
		},
	}
}