package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Healthcheck() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method == http.MethodGet && c.Path() == "/healthcheck" {
				return performHealthCheck(c)
			}
			return next(c)
		}
	}
}

func performHealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
