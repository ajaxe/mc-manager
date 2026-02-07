//go:build !wasm

package main

import (
	"context"
	"log"

	"github.com/ajaxe/mc-manager/internal/config"
	"github.com/ajaxe/mc-manager/internal/db"
	"github.com/ajaxe/mc-manager/internal/server"
	"github.com/labstack/echo/v4"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type empty struct{ app.Compo }

func Frontend() {
	app.RouteWithRegexp("/.*", app.NewZeroComponentFactory(&empty{}))
}

func Backend(ah *app.Handler) {
	cfg := config.LoadAppConfig()

	dbClient, err := db.NewClient(cfg)
	if err != nil {
		log.Fatalf("failed to create db client: %v", err)
	}
	defer func() {
		if err := dbClient.Close(context.Background()); err != nil {
			log.Printf("failed to close db client: %v", err)
		}
	}()

	s := server.NewBackendApi(dbClient)

	s.GET("/*", func(c echo.Context) error {
		ah.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	server.Start(s, dbClient)

	log.Println("end of Backend(...)")
}
