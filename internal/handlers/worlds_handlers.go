package handlers

import (
	"crypto/rand"
	"fmt"
	"os"
	"path"
	"regexp"
	"slices"
	"sort"
	"time"

	"github.com/ajaxe/mc-manager/internal/config"
	"github.com/ajaxe/mc-manager/internal/db"
	"github.com/ajaxe/mc-manager/internal/http"
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
	e.PUT("/worlds/:id", h.UpdateWorld("id"))
	e.DELETE("/worlds/:id", h.DeleteWorld("id"))
}

type worldsHandler struct {
	logger echo.Logger
}

func (w *worldsHandler) Worlds() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		d, err := db.Worlds()
		if err != nil {
			return
		}

		names, err := NewGameService(w.logger).serverDetails()
		if err != nil {
			return models.ErrAppGeneric(err)
		}

		for i := range d {
			d[i].IsActive = slices.ContainsFunc(names, func(s *models.GameServerDetail) bool {
				return toContainerName(d[i].Name) == s.Name
			})
		}

		sort.Slice(d, func(i, j int) bool {
			return d[i].IsActive == true && d[j].IsActive == false
		})

		return c.JSON(http.StatusOK, &models.WorldItemListResult{
			Data: d,
			ApiResult: models.ApiResult{
				Success: true,
			},
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

		if err := w.validateWorldName(u.Name); err != nil {
			return models.NewAppError(http.StatusBadRequest, err.Error(), err)
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

		return c.JSON(http.StatusOK, models.NewApiIDResult(id.Hex()))
	}
}
func (w *worldsHandler) DeleteWorld(idParam string) echo.HandlerFunc {
	return func(c echo.Context) error {
		i := c.Param(idParam)
		id, err := bson.ObjectIDFromHex(i)
		if err != nil {
			return models.ErrAppBadID(err)
		}

		if err := w.deleteWorld(id); err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}
func (w *worldsHandler) UpdateWorld(idParam string) echo.HandlerFunc {
	return func(c echo.Context) error {
		i := c.Param(idParam)
		id, err := bson.ObjectIDFromHex(i)
		if err != nil {
			return models.ErrAppBadID(err)
		}

		u := &models.WorldItem{}
		if err := c.Bind(u); err != nil {
			return models.NewAppError(http.StatusBadRequest, "Bad data.", nil)
		}

		if err := db.UpdateWorldByID(id, u); err != nil {
			return models.ErrAppGeneric(err)
		}

		return c.JSON(http.StatusOK, models.SuccessApiResult())
	}
}

func (w *worldsHandler) deleteWorld(id bson.ObjectID) (err error) {
	ww, err := db.WorldById(id)
	if err != nil {
		return
	}

	w.logger.Infof("found world by id: %v", ww.ID)

	if ww.IsFavorite {
		return models.NewAppError(http.StatusBadRequest, "Cannot delete world labeled as 'Favorite'", nil)
	}

	if err := db.DeleteWorldByID(id); err != nil {
		return err
	}

	cfg := config.LoadAppConfig()

	p := path.Join(cfg.GameServer.WorldDir, ww.Name)

	if _, e := os.Stat(p); e != nil {
		w.logger.Warnf("path: %s does not exist: %v", p, e)
	}

	if e := os.RemoveAll(p); e != nil {
		w.logger.Warnf("attempted to remove dir: %s failed: %v", p, e)
	}

	return
}
func (w *worldsHandler) validateWorldName(n string) error {
	expr := regexp.MustCompile(`^[a-zA-Z0-9 _-]{4,32}$`)

	if len(n) < 4 || len(n) > 32 {
		return fmt.Errorf("World name should be 4 to 32 characters long.")
	}

	if !expr.MatchString(n) {
		return fmt.Errorf("Invalid world name. Allowed character are: a to z, 'space', - and _")
	}

	return nil
}
