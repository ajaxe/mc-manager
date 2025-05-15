package pages

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type SetupPage struct {
	app.Compo
}

func (s *SetupPage) Render() app.UI {
	return &MainLayout{
		Content: []app.UI{
			app.H1().Text("Setup Page"),
			app.P().Text("This is the setup page."),
		},
	}
}
