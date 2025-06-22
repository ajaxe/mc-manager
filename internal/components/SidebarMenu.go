package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type SidebarMenu struct {
	app.Compo
	activeItem menuMap
	mapping    map[string]menuMap
}

func (s *SidebarMenu) OnNav(ctx app.Context) {
	p := app.Window().URL().Path
	if p == "/" {
		p = "/worlds"
	}
	s.activeItem = s.mapping[p]
}

func (s *SidebarMenu) Render() app.UI {
	if s.mapping == nil {
		s.initMapping()
	}

	return app.Div().
		Class("d-flex align-items-start flex-column bg-dark-subtle rounded-3 border-secondary-subtle sidebar").
		Style("height", "calc(100vh / 2)").
		Style("min-height", "400px").
		Body(
			app.Ul().Class("mb-auto p-2 nav flex-column").Body(
				s.navLink(s.mapping["/worlds"]),
				s.navLink(s.mapping["/worlds/add"]),
				s.navLink(s.mapping["/launches"]),
				s.navLink(s.mapping["/playtimer"]),
			),
			app.Ul().Class("p-2 nav flex-column").Body(
				s.navLink(s.mapping["/setup"]),
			),
		)
}

func (s *SidebarMenu) initMapping() {
	s.mapping = map[string]menuMap{
		"/worlds": {
			index: 0,
			name:  "Worlds",
			path:  "/worlds",
			icon:  "bi-globe-americas",
		},
		"/worlds/add": {
			index: 1,
			name:  "Add world",
			path:  "/worlds/add",
			icon:  "bi-plus-square",
		},
		"/launches": {
			index: 2,
			name:  "Launches",
			path:  "/launches",
			icon:  "bi-rocket-takeoff",
		},
		"/playtimer": {
			index: 2,
			name:  "Play Timer",
			path:  "/playtimer",
			icon:  "bi-hourglass-split",
		},
		"/setup": {
			index: 3,
			name:  "Setup",
			path:  "/setup",
			icon:  "bi-gear-fill",
		},
	}
}

func (s *SidebarMenu) navLink(item menuMap) app.UI {
	css := ""
	if item.path == s.activeItem.path {
		css = "text-bg-primary"
	}
	return app.Li().Class("nav-item").Body(
		app.A().
			Href(item.path).
			Class("nav-link rounded-3 "+css).
			Body(
				app.I().Class("bi "+item.icon),
				app.Div().Class("ms-3 d-none d-md-inline-block").Text(item.name),
			),
	)
}

type menuMap struct {
	index int
	name  string
	path  string
	icon  string
}
