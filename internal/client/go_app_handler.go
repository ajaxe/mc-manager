package client

import "github.com/maxence-charriere/go-app/v10/pkg/app"

var GoAppHandler = &app.Handler{
	Name:        "Minecraft Server Manager",
	Title:       "Minecraft Server Manager",
	Description: "Minecraft Server Manager",
	Icon:        app.Icon{Default: "/web/favicon.png", SVG: "/web/favicon.png"},

	Styles: []string{"/web/css/bootstrap.min.css", "/web/css/common.css", "/web/font/bootstrap-icons.min.css"},
	Scripts: []string{"/web/scripts/bootstrap.bundle.min.js",
		"/web/scripts/popper.min.js",
		"/web/scripts/cash.min.js",
		"/web/scripts/common.js",
	},
	Fonts: []string{"/web/font/fonts/bootstrap-icons.woff2"},
	HTML: func() app.HTMLHtml {
		return app.Html().DataSet("bs-theme", "dark")
	},
}
