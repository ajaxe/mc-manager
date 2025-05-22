package main

import (
	"github.com/ajaxe/mc-manager/internal/client"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func main() {
	Frontend()

	app.RunWhenOnBrowser()

	Backend(client.GoAppHandler)
}
