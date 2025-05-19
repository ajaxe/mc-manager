package components

import (
	"fmt"
	"time"

	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type LaunchItemCard struct {
	app.Compo
	Item *models.LaunchItem
}

func (l *LaunchItemCard) Render() app.UI {
	return app.Div().
		Class("card").
		Body(
			app.Div().
				Class("card-body").
				Body(
					app.Div().Class("row").Body(
						app.Div().
							Class("col").
							Text(l.Item.Name),
						app.Div().
							Class("col text-capitalize").
							Text(l.Item.GameMode),
						app.Div().
							Class("col-3").
							Text(l.Item.WorldSeed),
						app.Div().
							Class("col").
							Text(l.dtDisplay()),
						app.Div().
							Class("col").
							Body(l.statusIcon()),
					),
				),
		)
}
func (l *LaunchItemCard) statusIcon() app.UI {
	//l.Item.Status
	ico := "bi-check-circle-fill"
	co := "text-success"

	if l.Item.Status != "success" {
		ico = "bi-exclamation-octagon-fill"
		co = "text-danger"
	}
	return app.I().
		Class(fmt.Sprintf("bi %s %s", ico, co))
}
func (l *LaunchItemCard) dtDisplay() string {
	dt := l.Item.LaunchDate
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
	app.Logf("offset: %v", offset)

	dur, e := time.ParseDuration(fmt.Sprintf("%dm", -(offset)))
	if e != nil {
		app.Logf("time duration parse error: %v", e)
		return dt
	}

	return v.Add(dur).Format(time.DateTime)
}
