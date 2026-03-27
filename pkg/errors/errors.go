package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func Wrap(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

// Common errors
func BadRequest(msg string) *AppError {
	return New(http.StatusBadRequest, msg)
}

func Unauthorized(msg string) *AppError {
	return New(http.StatusUnauthorized, msg)
}

func Forbidden(msg string) *AppError {
	return New(http.StatusForbidden, msg)
}

func NotFound(msg string) *AppError {
	return New(http.StatusNotFound, msg)
}

func Conflict(msg string) *AppError {
	return New(http.StatusConflict, msg)
}

func Internal(msg string) *AppError {
	return New(http.StatusInternalServerError, msg)
}

func InternalWrap(msg string, err error) *AppError {
	return Wrap(http.StatusInternalServerError, msg, err)
}
