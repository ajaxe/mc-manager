package handlers

import (
	"github.com/ajaxe/mc-manager/internal/http"
	"github.com/labstack/echo/v4"
)

func Healthcheck() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			m := c.Request().Method
			p := c.Request().RequestURI
			if m == http.MethodGet && p == "/healthcheck" {
				return performHealthCheck(c)
			}
			return next(c)
		}
	}
}

func performHealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
