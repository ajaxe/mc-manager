package handlers

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/ajaxe/mc-manager/internal/db"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func AddWorldsHandlers(e *echo.Group, l echo.Logger) {
	h := &worldsHandler{
		logger: l,
	}

	e.GET("/worlds", h.Worlds())
	e.POST("/worlds", h.CreateWorld())
	e.DELETE("/worlds/:id", h.DeleteWorld("id"))
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

		n, err := gameServerIntance()
		if err != nil {
			return models.ErrAppGeneric(err)
		}

		for i := range w {
			w[i].IsActive = w[i].Name == n
		}

		sort.Slice(w, func(i, j int) bool {
			return w[i].IsActive == true && w[j].IsActive == false
		})

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
func (w *worldsHandler) DeleteWorld(idParam string) echo.HandlerFunc {
	return func(c echo.Context) error {
		i := c.Param(idParam)
		id, err := bson.ObjectIDFromHex(i)
		if err != nil {
			return models.ErrAppBadID(err)
		}

		if err := db.DeleteWorldByID(id); err != nil {
			return models.ErrAppGeneric(err)
		}

		return c.NoContent(http.StatusNoContent)
	}
}
