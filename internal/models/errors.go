package models

import (
	"fmt"
	"strings"
)

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
