package client

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func fetch(url string, input *map[string]any) (json *string, code int, err error) {
	if url == "" {
		err = fmt.Errorf("URL not specified for Fetch()")
		return
	}

	chCode := make(chan int, 1)
	chJson := make(chan string, 1)
	chErr := make(chan error, 1)

	defer func(ch chan int) { close(ch) }(chCode)
	defer func(ch chan string) { close(ch) }(chJson)
	defer func(ch chan error) { close(ch) }(chErr)

	//app.Logf("invoking fetch: url: %s, input: %v", url, input)

	p1 := app.Window().Call("fetch", url, *input)
	p1.Call("then", app.FuncOf(func(this app.Value, args []app.Value) any {
		resp := args[0]
		p2 := resp.Call("json")
		p2.Call("then", app.FuncOf(func(this app.Value, args []app.Value) any {
			res := args[0]

			chCode <- resp.Get("status").Int()
			chJson <- app.Window().Get("JSON").Call("stringify", res).String()
			chErr <- nil

			return nil
		})).Call("catch", app.FuncOf(func(this app.Value, args []app.Value) any {
			chCode <- 500
			chErr <- fmt.Errorf("%s", args[0].Get("message"))
			chJson <- ""

			return nil
		}))

		return nil
	})).Call("catch", app.FuncOf(func(this app.Value, args []app.Value) any {
		chErr <- fmt.Errorf("%s", args[0].Get("message"))
		chJson <- ""

		return nil
	}))

	j := <-chJson
	json = &j
	code = <-chCode
	err = <-chErr

	return
}
