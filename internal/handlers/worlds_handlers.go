package handlers

import (
	"net/http"

	"github.com/ajaxe/mc-manager/internal/db"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/labstack/echo/v4"
)

func AddWorldsHandlers(e *echo.Group, l echo.Logger) {
	h := &worldsHandler{
		logger: l,
	}

	e.GET("/worlds", h.Worlds())
}

type worldsHandler struct {
	logger echo.Logger
}

func (w *worldsHandler) Worlds() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		/* w := []*models.WorldItem{
			{
				ID:          bson.NewObjectID(),
				Name:        "World 1",
				Description: "Description for World 1",
				WorldSeed:   "Seed 1",
			},
			{
				ID:          bson.NewObjectID(),
				Name:        "World 2",
				Description: "Description for World 2",
				WorldSeed:   "Seed 2",
			},
		} */
		w, err := db.Worlds()
		if err != nil {
			return
		}
		return c.JSON(http.StatusOK, &models.WorldItemListResult{
			Data: w,
		})
	}
}
