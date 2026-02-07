package handlers

import (
	"github.com/ajaxe/mc-manager/internal/db"
	"github.com/ajaxe/mc-manager/internal/http"
	"github.com/labstack/echo/v4"
)

func Healthcheck(db *db.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			m := c.Request().Method
			p := c.Request().RequestURI
			if m == http.MethodGet && p == "/healthcheck" {
				return performHealthCheck(c, db)
			}
			return next(c)
		}
	}
}

func performHealthCheck(c echo.Context, client *db.Client) error {
	if client == nil {
		return c.String(http.StatusInternalServerError, "DB client not initialized")
	}
	if err := client.Ping(); err != nil {
		return c.String(http.StatusInternalServerError, "DB ping failed")
	}
	return c.String(http.StatusOK, "OK")
}
