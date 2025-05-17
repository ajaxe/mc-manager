package models

import (
	"fmt"
	"net/http"
	"strings"
)

func ErrAppBadID(err error) *AppError {
	return NewAppError(http.StatusBadRequest, "Invalid ID.", err)
}

func ErrAppGeneric(err error) *AppError {
	return NewAppError(http.StatusInternalServerError, "Something went wrong.", err)
}
func ErrInvalidData(err error) *AppError {
	return NewAppError(http.StatusBadRequest, "Invalid data.", err)
}
func ErrAppRequired(m string) *AppError {
	return NewAppError(http.StatusBadRequest, fmt.Sprintf("%s is required.", m), nil)
}

func NewAppError(status int, message string, e error) *AppError {
	return &AppError{
		status:  status,
		message: message,
		err:     e,
	}
}

type AppError struct {
	status  int
	message string
	err     error
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %v", strings.ToLower(e.message), e.err)
}
func (e *AppError) HTTPStatus() int {
	return e.status
}
func (e *AppError) Message() string {
	return e.message
}
func (e *AppError) Inner() error {
	return e.err
}
