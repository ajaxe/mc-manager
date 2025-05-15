package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type SidebarMenu struct {
	app.Compo
}

func (s *SidebarMenu) Render() app.UI {
	return app.Div().
		Class("d-flex align-items-start flex-column bg-dark-subtle rounded border-secondary-subtle sidebar").
		Style("height", "calc(100vh / 2)").
		Style("min-height", "400px").
		Body(
			app.Ul().Class("mb-auto p-2 nav flex-column").Body(
				s.navLink("Worlds", "/worlds", "bi-globe-americas"),
				s.navLink("Add world", "/add", "bi-plus-square"),
				s.navLink("Launches", "/launches", "bi-rocket-takeoff"),
			),
			app.Ul().Class("p-2 nav flex-column").Body(
				s.navLink("Setup", "/setup", "bi-gear-fill"),
			),
		)
}
func (s *SidebarMenu) navLink(text, href, icon string) app.UI {
	return app.Li().Class("nav-item rounded").Body(
		app.A().
			Href(href).
			Class("nav-link").
			Body(
				app.I().Class("bi "+icon),
				app.Div().Class("d-inline-block ms-3").Text(text),
			),
	)
}
