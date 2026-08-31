package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ajaxe/mc-manager/internal/config"
	"github.com/ajaxe/mc-manager/internal/db"
	"github.com/ajaxe/mc-manager/internal/handlers"
	"github.com/ajaxe/mc-manager/internal/job"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func NewBackendApi(db *db.Client) *echo.Echo {
	// Setup slog to use JSON handler
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	e := echo.New()
	e.Logger = logger
	e.HTTPErrorHandler = handlers.AppErrorHandler()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			cc := models.NewAppContext(c, logger)
			c.Set("appContext", cc)
			return next(c)
		}
	})

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogLatency:   true,
		LogMethod:    true,
		LogUserAgent: true,
		LogRemoteIP:  true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			level := slog.LevelInfo
			if v.Error != nil || v.Status >= 500 {
				level = slog.LevelError
			} else if v.Status >= 400 {
				level = slog.LevelWarn
			}

			attrs := []any{
				"uri", v.URI,
				"status", v.Status,
				"latency_ms", v.Latency.Milliseconds(),
				"method", v.Method,
				"remote_ip", v.RemoteIP,
			}
			if v.UserAgent != "" {
				attrs = append(attrs, "user_agent", v.UserAgent)
			}
			if v.Error != nil {
				attrs = append(attrs, "error", v.Error.Error())
			}

			logger.Log(c.Request().Context(), level, "HTTP request", attrs...)
			return nil
		},
	}))
	e.Use(middleware.Recover())
	e.Use(handlers.Healthcheck(db))

	a := e.Group("/api")

	handlers.AddLoginHandlers(a, logger)
	handlers.AddWorldsHandlers(a, logger, db)
	handlers.AddLaunchHandlers(a, logger, db)
	handlers.AddPlaytimerHandlers(a, logger, db)

	return e
}

// Start echo server with graceful handling of process termination.
func Start(e *echo.Echo, db *db.Client) {
	cfg := config.LoadAppConfig()
	addr := fmt.Sprintf(":%v", cfg.Server.Port)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	defer stop()

	go func() {
		job.StartMonitor(ctx, e.Logger, db)
	}()

	sc := echo.StartConfig{
		Address:         addr,
		GracefulTimeout: 10 * time.Second,
	}

	var err error
	if cfg.UseTLS() {
		e.Logger.Info("starting server with tls")
		err = sc.StartTLS(ctx, e, cfg.Server.CertFile, cfg.Server.KeyFile)
	} else {
		e.Logger.Info("starting server without tls")
		err = sc.Start(ctx, e)
	}

	if err != nil && err != http.ErrServerClosed {
		e.Logger.Error("shutting down the server", "error", err)
		os.Exit(1)
	}

	e.Logger.Info("server exited properly")
}
