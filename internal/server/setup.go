package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log/slog"

	"github.com/ajaxe/mc-manager/internal/config"
	"github.com/ajaxe/mc-manager/internal/db"
	"github.com/ajaxe/mc-manager/internal/handlers"
	"github.com/ajaxe/mc-manager/internal/job"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	elog "github.com/labstack/gommon/log"
)

func NewBackendApi(db *db.Client) *echo.Echo {
	// Setup slog to use JSON handler
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	e := echo.New()
	e.Logger.SetLevel(elog.DEBUG)
	e.HTTPErrorHandler = handlers.AppErrorHandler()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := models.NewAppContext(c, e.Logger)
			return next(cc)
		}
	})

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(handlers.Healthcheck(db))

	a := e.Group("/api")

	handlers.AddLoginHandlers(a, e.Logger)
	handlers.AddWorldsHandlers(a, e.Logger, db)
	handlers.AddLaunchHandlers(a, e.Logger, db)
	handlers.AddPlaytimerHandlers(a, e.Logger, db)

	return e
}

// Start echo server with graceful hanlding of process termination.
func Start(e *echo.Echo, db *db.Client) {
	cfg := config.LoadAppConfig()
	addr := fmt.Sprintf(":%v", cfg.Server.Port)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	defer stop()
	// Start server
	go func() {
		var err error
		if cfg.UseTLS() {
			e.Logger.Info("starting server with tls")
			err = e.StartTLS(addr, cfg.Server.CertFile, cfg.Server.KeyFile)
		} else {
			e.Logger.Info("starting server without tls")
			err = e.Start(addr)
		}
		if err != nil && err != http.ErrServerClosed {
			e.Logger.Fatalf("shutting down the server: %v", err)
		}
	}()

	go func() {
		job.StartMonitor(ctx, e.Logger, db)
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("failed to shutdown server: %v", err)
	}
	e.Logger.Info("server exited properly")
}
