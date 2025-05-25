package handlers

import (
	"github.com/ajaxe/mc-manager/internal/http"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/labstack/echo/v4"
)

func AddLoginHandlers(e *echo.Group, l echo.Logger) {
	h := &loginHandler{
		logger: l,
	}

	e.POST("/login/check", h.check())
}

type loginHandler struct {
	logger echo.Logger
}

func (l *loginHandler) check() echo.HandlerFunc {
	return func(c echo.Context) error {
		/* return c.JSON(http.StatusOK, models.ApiResult{
			Success: true,
		}) */
		return c.JSON(http.StatusOK, &models.ApiResult{
			Success: true,
		})
	}
}
