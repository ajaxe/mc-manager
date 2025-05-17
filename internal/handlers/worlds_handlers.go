package handlers

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"

	"github.com/ajaxe/mc-manager/internal/db"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/labstack/echo/v4"
)

func AddWorldsHandlers(e *echo.Group, l echo.Logger) {
	h := &worldsHandler{
		logger: l,
	}

	e.GET("/worlds", h.Worlds())
	e.POST("/worlds", h.CreateWorld())
}

type worldsHandler struct {
	logger echo.Logger
}

func (w *worldsHandler) Worlds() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		w, err := db.Worlds()
		if err != nil {
			return
		}
		return c.JSON(http.StatusOK, &models.WorldItemListResult{
			Data: w,
		})
	}
}

func (w *worldsHandler) CreateWorld() echo.HandlerFunc {
	return func(c echo.Context) error {
		u := &models.WorldItem{}
		if err := c.Bind(u); err != nil {
			return models.NewAppError(http.StatusBadRequest, "Bad data.", nil)
		}

		if u.Name == "" {
			return models.NewAppError(http.StatusBadRequest, "World name is required.", nil)
		}

		if u.WorldSeed == "" {
			u.WorldSeed = rand.Text()
		}

		if u.Description == "" {
			u.Description = u.Name
		}

		if u.GameMode == "" {
			u.GameMode = "survival"
		}

		u.CreateDate = time.Now().UTC().Format(time.RFC3339)
		id, err := db.InsertWorld(u)
		if err != nil {
			return models.ErrAppGeneric(fmt.Errorf("error saving user: %v", err))
		}

		return c.JSON(http.StatusOK, models.NewApiIDResult(id))
	}
}
