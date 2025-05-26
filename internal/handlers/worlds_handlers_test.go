package handlers

import (
	"testing"

	"github.com/labstack/gommon/log"
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

func createWorldHandlerSut() *worldsHandler {
	return &worldsHandler{
		logger: log.New("echo_test"),
	}
}
