package server

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func TestNewBackendApi_JSONRequestLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	e := echo.New()
	e.Logger = logger

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

	e.GET("/test", func(c *echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("User-Agent", "Go-Test-Agent")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	logStr := buf.String()
	if logStr == "" {
		t.Fatal("expected non-empty log output")
	}

	var entry map[string]any
	if err := json.Unmarshal([]byte(logStr), &entry); err != nil {
		t.Fatalf("expected valid JSON log entry, got error: %v, log: %s", err, logStr)
	}

	if entry["msg"] != "HTTP request" {
		t.Errorf("expected msg 'HTTP request', got %v", entry["msg"])
	}

	if entry["uri"] != "/test" {
		t.Errorf("expected uri '/test', got %v", entry["uri"])
	}

	if entry["method"] != "GET" {
		t.Errorf("expected method 'GET', got %v", entry["method"])
	}

	if entry["user_agent"] != "Go-Test-Agent" {
		t.Errorf("expected user_agent 'Go-Test-Agent', got %v", entry["user_agent"])
	}
}
