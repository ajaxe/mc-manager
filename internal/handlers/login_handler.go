package handlers

import (
	"github.com/ajaxe/mc-manager/internal/config"
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
	return func(c echo.Context) (err error) {
		cfg := config.LoadAppConfig()
		if cfg.Server.AuthRedirectURL == "" {
			return c.JSON(http.StatusOK, models.NewApiAuthResult())
		}

		u, err := cfg.AuthRedirectURL()
		if err != nil {
			l.logger.Errorf("failed to read auth redirect url: %v", err)
			return
		}

		ck := c.Request().Cookies()
		for _, c := range ck {
			if c.Name == cfg.Server.AuthCookieName {
				u = ""
				l.logger.Infof("authenticated session detected, no need to redirect.")
				break
			}
		}

		if u != "" {
			l.logger.Infof("un-authenticated session detected, will redirect to [%s].", u)
		}

		return c.JSON(http.StatusOK, models.NewApiAuthResult(u))
	}
}
