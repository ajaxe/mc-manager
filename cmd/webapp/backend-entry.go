//go:build !wasm

package main

import (
	"log"

	"github.com/ajaxe/mc-manager/internal/server"
	"github.com/labstack/echo/v4"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type empty struct{ app.Compo }

func Frontend() {
	app.RouteWithRegexp("/.*", app.NewZeroComponentFactory(&empty{}))
}

func Backend(ah *app.Handler) {
	s := server.NewBackendApi()

	s.GET("/*", func(c echo.Context) error {
		ah.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	server.Start(s)

	log.Println("end of Backend(...)")
}
