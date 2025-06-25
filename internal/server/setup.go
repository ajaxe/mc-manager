package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ajaxe/mc-manager/internal/config"
	"github.com/ajaxe/mc-manager/internal/handlers"
	"github.com/ajaxe/mc-manager/internal/job"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	elog "github.com/labstack/gommon/log"
)

func NewBackendApi() *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(elog.DEBUG)
	e.HTTPErrorHandler = handlers.AppErrorHandler()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(handlers.Healthcheck())

	a := e.Group("/api")

	handlers.AddLoginHandlers(a, e.Logger)
	handlers.AddWorldsHandlers(a, e.Logger)
	handlers.AddLaunchHandlers(a, e.Logger)
	handlers.AddPlaytimerHandlers(a, e.Logger)

	return e
}

// Start echo server with graceful hanlding of process termination.
func Start(e *echo.Echo) {
	cfg := config.LoadAppConfig()
	addr := fmt.Sprintf(":%v", cfg.Server.Port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
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
		job.StartMonitor(ctx, e.Logger)
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("failed to shutdown server: %v", err)
	}
}
