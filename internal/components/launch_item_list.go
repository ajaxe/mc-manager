package components

import (
	"github.com/ajaxe/mc-manager/internal/client"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type LaunchItemList struct {
	app.Compo
	items     models.LaunchItemListResult
	incoming  models.LaunchItemListResult
	direction string
}

func (l *LaunchItemList) OnMount(ctx app.Context) {
	ctx.ObserveState(client.StateKeyLaunches, &l.incoming)
}
func (l *LaunchItemList) Render() app.UI {
	if len(l.incoming.Results) != 0 {
		l.items = l.incoming
	}

	l.items.PrevID = l.incoming.PrevID
	l.items.NextID = l.incoming.NextID

	if !l.incoming.HasMore {
		if l.direction == models.PageDirectionPrev {
			l.items.PrevID = ""
		}
		if l.direction == models.PageDirectionNext {
			l.items.NextID = ""
		}
	}

	r := l.items.Results

	return app.Div().
		//Class("row g-4").
		Body(
			&TablePager{
				NextID: l.items.NextID,
				PrevID: l.items.PrevID,
				OnClick: func(ctx app.Context, s string) {
					id := l.items.PrevID
					if s == models.PageDirectionNext {
						id = l.items.NextID
					}
					l.direction = s
					client.NewAppContext(ctx).
						LoadData(client.StateKeyLaunches, models.LaunchItemListRequest{
							Direction: s,
							CursorID:  id,
						})
				},
			},
			l.header(),
			app.Range(r).
				Slice(func(i int) app.UI {
					return &LaunchItemCard{
						Item: r[i],
					}
				}),
		)
}

func (l *LaunchItemList) header() app.UI {
	return app.Div().
		Class("card d-none d-md-block").
		Body(
			app.Div().
				Class("card-body fw-bold").
				Body(
					app.Div().Class("row").Body(
						app.Div().
							Class("col").
							Text("World Name"),
						app.Div().
							Class("col text-capitalize").
							Text("GameMode"),
						app.Div().
							Class("col-3").
							Text("World Seed"),
						app.Div().
							Class("col").
							Text("Launch Date"),
						app.Div().
							Class("col").
							Text("Status"),
					),
				),
		)
}
