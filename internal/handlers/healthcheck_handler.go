package handlers

import (
	"github.com/ajaxe/mc-manager/internal/db"
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
	_, err := db.NewClient()
	if err != nil {
		return c.String(http.StatusInternalServerError, "DB connection failed")
	}
	err = db.Ping()
	if err != nil {
		return c.String(http.StatusInternalServerError, "DB ping failed")
	}
	return c.String(http.StatusOK, "OK")
}
