package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/ajaxe/mc-manager/internal/db"
	"github.com/ajaxe/mc-manager/internal/http"
	"github.com/ajaxe/mc-manager/internal/job"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func AddPlaytimerHandlers(e *echo.Group, l echo.Logger, db *db.Client) {
	h := &playTimerHandler{
		logger:      l,
		gameService: NewGameService(l),
		db:          db,
	}

	e.GET("/playtimer", h.PlayTimer())
	e.POST("/playtimer", h.CreatePlayTimer())
	e.DELETE("/playtimer", h.DeletePlayTimer())
}

type playTimerHandler struct {
	logger      echo.Logger
	gameService GameService
	db          *db.Client
}

// PlayTimer returns a GET handler which responds with active play timer item.
func (p *playTimerHandler) PlayTimer() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		pt, err := p.db.ActivePlayTimer(c.Request().Context())
		if err != nil {
			err = models.ErrAppGeneric(err)
			return
		}
		d := []*models.PlayTimerItem{}
		if pt != nil {
			d = append(d, pt)
		}
		return c.JSON(http.StatusOK, &models.PlayTimerListResult{
			Data: d,
			ApiResult: models.ApiResult{
				Success: true,
			},
		})
	}
}

// CreatePlayTimer returns a POST handler which creates a new play timer item.
func (p *playTimerHandler) CreatePlayTimer() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		cc := c.(*models.AppContext)

		if cc.IsAdmin() == false {
			return models.NewAppError(http.StatusBadRequest, "Only Poacha can create play timer.", nil)
		}

		u := &models.PlayTimerItem{}
		if err := c.Bind(u); err != nil {
			return models.NewAppError(http.StatusBadRequest, "Bad data.", nil)
		}

		id, err := p.create(c.Request().Context(), u)
		if err != nil {
			return
		}

		return cc.JSON(http.StatusOK, models.NewApiIDResult(id.Hex()))
	}
}

// DeletePlayTimer returns a DELETE handler which deactivates the current play timer item.
func (p *playTimerHandler) DeletePlayTimer() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		cc := c.(*models.AppContext)

		if cc.IsAdmin() == false {
			return models.NewAppError(http.StatusBadRequest, "Only Poacha can delete play timer.", nil)
		}

		pt, err := p.db.ActivePlayTimer(c.Request().Context())
		if pt != nil {
			id, _ := bson.ObjectIDFromHex(pt.ID)
			pt.IsActive = false
			err = p.db.UpdatePlayTimerByID(c.Request().Context(), id, pt)
			if err != nil {
				err = models.ErrAppGeneric(err)
				return
			}
			job.StopCurrentPlayTimer()
		}
		return cc.JSON(http.StatusOK, models.SuccessApiResult())
	}
}

// create inserts a new play timer item into the database.
// It also deactivates any currently active play timer item.
func (p *playTimerHandler) create(ctx context.Context, d *models.PlayTimerItem) (id bson.ObjectID, err error) {
	if d.Minutes <= 0 {
		err = models.NewAppError(http.StatusBadRequest, "Invalid minutes value.", nil)
		return
	}

	now := time.Now().UTC()
	d.IsActive = true
	d.CreateDate = now.Format(time.RFC3339)
	d.EndDate = now.Add(time.Minute * time.Duration(d.Minutes)).Format(time.RFC3339)

	active, err := p.db.ActivePlayTimer(ctx)
	if err != nil {
		err = models.ErrAppGeneric(err)
		return
	}

	if active != nil {
		active.IsActive = false
		i, _ := bson.ObjectIDFromHex(active.ID)
		err = p.db.UpdatePlayTimerByID(ctx, i, active)
		if err != nil {
			err = models.ErrAppGeneric(err)
			return
		}
	}

	id, err = p.db.InsertPlayTimer(ctx, d)
	if err != nil {
		err = models.ErrAppGeneric(err)
	}

	if err == nil {
		e := p.gameService.sendMessageToServer(fmt.Sprintf("A new play timer has been set for %d minutes.", d.Minutes))
		if e != nil {
			p.logger.Warnf("Failed to send message to server: %v", e)
		}
	}

	job.QueueJob(d)

	return
}
