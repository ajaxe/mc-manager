package handlers

import (
	"fmt"

	"github.com/ajaxe/mc-manager/internal/http"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/labstack/echo/v5"
)

func AppErrorHandler() echo.HTTPErrorHandler {
	return func(c *echo.Context, err error) {
		if rw, _ := echo.UnwrapResponse(c.Response()); rw != nil && rw.Committed {
			return
		}
		c.Logger().Error(fmt.Sprintf("error: %v", err))
		st := http.StatusInternalServerError
		m := err.Error()
		if ae, ok := err.(*models.AppError); ok {
			st = ae.HTTPStatus()
			m = ae.Message()
		}
		c.JSON(st, &models.ApiResult{
			Success:      false,
			ErrorMessage: m,
		})
	}
}
