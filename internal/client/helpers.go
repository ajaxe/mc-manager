package client

import (
	"fmt"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)
// BrowserDateDisplay converts a date string in RFC3339 format to a browser-friendly display format.
func BrowserDateDisplay(dt string) string {
	v, e := time.Parse(time.RFC3339, dt)

	if e != nil {
		app.Logf("time parse error: %v", e)
		return dt
	}

	offset := app.Window().
		Get("Date").
		New().
		Call("getTimezoneOffset").
		Int()

	dur, e := time.ParseDuration(fmt.Sprintf("%dm", -(offset)))
	if e != nil {
		app.Logf("time duration parse error: %v", e)
		return dt
	}

	return v.Add(dur).Format("Mon, Jan 2 2006 3:04 PM")
}
