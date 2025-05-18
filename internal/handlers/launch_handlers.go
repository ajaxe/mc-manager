package handlers

import (
	"fmt"
	"net/http"

	"github.com/ajaxe/mc-manager/internal/db"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func AddLaunchHandlers(e *echo.Group, l echo.Logger) {
	h := &launchHandler{
		logger: l,
	}

	e.GET("/launches", h.Launches())
	e.POST("/launches", h.CreateLaunch())
	e.DELETE("/launches/:id", h.DeleteLaunch("id"))
}

type launchHandler struct {
	logger echo.Logger
}

func (l *launchHandler) DeleteLaunch(s string) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		return c.String(http.StatusOK, "ok")
	}
}

func (l *launchHandler) CreateLaunch() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		u := &models.CreateLaunchItem{}
		if err := c.Bind(u); err != nil {
			return models.NewAppError(http.StatusBadRequest, "Bad data.", nil)
		}

		if u.WorldItemID.String() == "" || u.WorldItemID == bson.NilObjectID {
			return models.NewAppError(http.StatusBadRequest, "Bad data.", nil)
		}

		w, err := db.WorlById(u.WorldItemID)
		if err != nil {
			return models.ErrAppGeneric(fmt.Errorf("world not found: %v", err))
		}

		existing, err := gameServerIntance()
		if err != nil {
			return models.ErrAppGeneric(err)
		}

		if existing != "" {
			l.logger.Info("Game server already running, cannot create new launch.")
			return c.String(http.StatusOK, "ok")
		}

		if err := createGameServer(w); err != nil {
			l.logger.Error("Failed to create game server: %v", err)
			return models.ErrAppGeneric(err)
		}

		return c.JSON(http.StatusOK, &models.ApiResult{
			Success: true,
		})
	}
}

func (l *launchHandler) Launches() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		l, err := db.Launches()
		if err != nil {
			return
		}

		return c.JSON(http.StatusOK, &models.LaunchItemListResult{
			Data: l,
			ApiResult: models.ApiResult{
				Success: true,
			},
		})
	}
}
