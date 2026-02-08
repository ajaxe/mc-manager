package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ajaxe/mc-manager/internal/db"
	"github.com/ajaxe/mc-manager/internal/http"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func AddLaunchHandlers(e *echo.Group, l echo.Logger, db *db.Client) {
	h := &launchHandler{
		logger:      l,
		gameService: NewGameService(l),
		db:          db,
	}

	e.GET("/launches", h.Launches())
	e.POST("/launches", h.CreateLaunch())
	e.DELETE("/launches/:id", h.DeleteLaunch("id"))
}

type launchHandler struct {
	logger      echo.Logger
	gameService GameService
	db          *db.Client
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

		if u.WorldItemID == "" || u.WorldItemID == bson.NilObjectID.Hex() {
			return models.NewAppError(http.StatusBadRequest, "Bad data.", nil)
		}

		id, err := l.createLaunch(c.Request().Context(), u)
		if err != nil {
			return
		}

		l.logger.Info("Created new launch item: %s", id.Hex())

		return c.JSON(http.StatusOK, models.NewApiIDResult(id.Hex()))
	}
}

func (l *launchHandler) Launches() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		dir := c.QueryParam("dir")
		curID := c.QueryParam("cursorId")
		pgs := c.QueryParam("pageSize")

		paged, err := l.list(c.Request().Context(), dir, curID, pgs)

		return c.JSON(http.StatusOK, &models.LaunchItemListResult{
			PaginationResult: paged,
			ApiResult: models.ApiResult{
				Success: true,
			},
		})
	}
}

func (l *launchHandler) list(ctx context.Context, dir, curID, pgs string) (paged models.PaginationResult[models.LaunchItem], err error) {
	pg := 10
	if n, _ := strconv.Atoi(pgs); pgs != "" {
		pg = n
	} else {
		l.logger.Info("using default size 10")
	}

	if dir == "" {
		dir = models.PageDirectionNext
	}

	paged, err = l.db.Launches(ctx, db.PaginationOptions{
		Direction: dir,
		PageSize:  pg,
		CursorID:  curID,
	})

	return
}

func (l *launchHandler) createLaunch(ctx context.Context, u *models.CreateLaunchItem) (id bson.ObjectID, err error) {
	i, err := bson.ObjectIDFromHex(u.WorldItemID)
	if err != nil {
		return
	}
	w, err := l.db.WorldById(ctx, i)
	if err != nil {
		err = models.ErrAppGeneric(fmt.Errorf("world not found: %v", err))
		return
	}

	gameService := l.gameService

	existing, err := gameService.serverDetails()
	if err != nil {
		err = models.ErrAppGeneric(err)
		return
	}

	if l.checkLauchItem(existing, w) {
		l.logger.Info("Game server already running, cannot create new launch.")
		err = models.ErrAppGeneric(fmt.Errorf("Game server already running: %s", w.Name))
		return
	}

	err = gameService.stopAllInstances()
	if err != nil {
		err = models.ErrAppGeneric(err)
		return
	}

	if err = gameService.createGameServer(w); err != nil {
		l.logger.Error("Failed to create game server: %v", err)

		_, e := l.db.LaunchInsert(ctx, models.ToLaunchItem(w, time.Now().UTC().Format(time.RFC3339), "failed"))

		l.logger.Error("Failed to insert launch item: %v", e)

		err = models.ErrAppGeneric(err)
		return
	}

	id, err = l.db.LaunchInsert(ctx, models.ToLaunchItem(w, time.Now().UTC().Format(time.RFC3339), "success"))

	if err != nil {
		err = models.ErrAppGeneric(err)
		//TODO: remove started container
		return
	}

	return
}

// checks if the "world name" is already in the list of running game servers "names".
// returns true if the world name is found in the list
func (l *launchHandler) checkLauchItem(details []*models.GameServerDetail, world *models.WorldItem) bool {
	for _, d := range details {
		if d.Name == toContainerName(world.Name) && d.GameMode == world.GameMode {
			return true
		}
	}
	return false
}
