package handlers

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/labstack/echo/v5"
)

func TestValidateWorldName_InvalidLength(t *testing.T) {
	sut := createWorldHandlerSut()
	expected := "World name should be 4 to 32 characters long."
	test := "123"
	if e := sut.validateWorldName(test); e == nil || e.Error() != expected {
		t.Fatalf("expect validation error: '%s', got: %v", expected, e)
	}
}

func TestValidateWorldName_InvalidChars(t *testing.T) {
	sut := createWorldHandlerSut()
	expected := "Invalid world name. Allowed character are: a to z, 'space', - and _"
	test := "123 # @ _"
	if e := sut.validateWorldName(test); e == nil || e.Error() != expected {
		t.Fatalf("expect validation error: '%s', got: %v", expected, e)
	}
}

func TestValidateWorldName_Valid(t *testing.T) {
	sut := createWorldHandlerSut()

	test := "Test world-12_3"
	if e := sut.validateWorldName(test); e != nil {
		t.Fatalf("do not expect error. got: %v", e)
	}
}

func TestAppErrorHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := AppErrorHandler()
	appErr := models.NewAppError(http.StatusBadRequest, "Invalid test input", nil)

	handler(c, appErr)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func createWorldHandlerSut() *worldsHandler {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	return &worldsHandler{
		logger: logger,
	}
}
