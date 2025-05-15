package main

import (
	"github.com/ajaxe/mc-manager/internal/client"
	"github.com/ajaxe/mc-manager/internal/pages"
	"github.com/ajaxe/mc-manager/internal/server"
	"github.com/labstack/echo/v4"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func main() {
	// This is the entry point for the web application.
	// The main function will initialize the application and start the server.
	app.Route("/", func() app.Composer { return &pages.HomePage{} })
	app.Route("/worlds", func() app.Composer { return &pages.HomePage{} })

	app.RunWhenOnBrowser()

	s := server.NewBackendApi()

	s.GET("/*", func(c echo.Context) error {
		client.GoAppHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	server.Start(s)
}
