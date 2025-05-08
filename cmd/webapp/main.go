package main

import (
	"log"
	"net/http"

	"github.com/ajaxe/mc-manager/internal/pages"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func main() {
	// This is the entry point for the web application.
	// The main function will initialize the application and start the server.
	app.Route("/", func() app.Composer { return &pages.HomePage{} })

	app.RunWhenOnBrowser()

	http.Handle("/", &app.Handler{
		Name:        "Hello",
		Description: "An Hello World! example",
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
