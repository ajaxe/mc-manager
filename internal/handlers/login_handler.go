package handlers

import (
	"fmt"
	"log/slog"

	"github.com/ajaxe/mc-manager/internal/config"
	"github.com/ajaxe/mc-manager/internal/http"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/labstack/echo/v5"
)

func AddLoginHandlers(e *echo.Group, l *slog.Logger) {
	h := &loginHandler{
		logger: l,
	}

	e.POST("/login/check", h.check())
}

type loginHandler struct {
	logger *slog.Logger
}

func (l *loginHandler) check() echo.HandlerFunc {
	return func(c *echo.Context) (err error) {
		cc := models.GetAppContext(c)
		cfg := config.LoadAppConfig()
		if cfg.Server.AuthServerEnabled == false {
			return cc.JSON(http.StatusOK, models.NewApiAuthResult())
		}

		u, err := cfg.AuthRedirectURL()
		if err != nil {
			l.logger.Error(fmt.Sprintf("failed to read auth redirect url: %v", err))
			return
		}

		ck := cc.AuthCookieValue()
		if ck != "" {
			u = ""
			l.logger.Info("authenticated session detected, no need to redirect.")
		}

		if u != "" {
			l.logger.Info(fmt.Sprintf("un-authenticated session detected, will redirect to [%s].", u))
		}

		return cc.JSON(http.StatusOK, models.NewApiAuthResult(u))
	}
}
