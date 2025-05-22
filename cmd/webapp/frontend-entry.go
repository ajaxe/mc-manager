//go:build wasm

package main

import (
	"github.com/ajaxe/mc-manager/internal/pages"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func Frontend() {
	// This is the entry point for the web application.
	// The main function will initialize the application and start the server.
	app.Route("/", func() app.Composer { return &pages.HomePage{} })
	app.Route("/worlds", func() app.Composer { return &pages.HomePage{} })
	app.Route("/worlds/add", func() app.Composer { return &pages.WorldAddPage{} })
	app.Route("/launches", func() app.Composer { return &pages.WorldLaunchPage{} })
	app.Route("/setup", func() app.Composer { return &pages.SetupPage{} })
}

func Backend(_ *app.Handler) {

}
