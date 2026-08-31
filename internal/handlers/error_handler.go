package handlers

import (
	"github.com/ajaxe/mc-manager/internal/http"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/labstack/echo/v5"
)

func AppErrorHandler() echo.HTTPErrorHandler {
	return func(c *echo.Context, err error) {
		if rw, _ := echo.UnwrapResponse(c.Response()); rw != nil && rw.Committed {
			return
		}
		st := http.StatusInternalServerError
		m := err.Error()
		if ae, ok := err.(*models.AppError); ok {
			st = ae.HTTPStatus()
			m = ae.Message()
		}
		c.Logger().Error("HTTP request handler error", "error", err, "status", st)
		c.JSON(st, &models.ApiResult{
			Success:      false,
			ErrorMessage: m,
		})
	}
}
