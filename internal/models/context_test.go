package models

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
)

func TestGetAppContext(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	appCtx := NewAppContext(c, logger)
	c.Set("appContext", appCtx)

	retrieved := GetAppContext(c)
	if retrieved == nil {
		t.Fatal("expected non-nil AppContext")
	}
	if retrieved != appCtx {
		t.Fatal("expected retrieved AppContext to match stored instance")
	}
}

func TestGetAppContext_DefaultFallback(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	retrieved := GetAppContext(c)
	if retrieved == nil {
		t.Fatal("expected non-nil AppContext fallback")
	}
}
