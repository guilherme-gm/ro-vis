package controllers

import (
	"fmt"
	"net/http"
)

type HttpError struct {
	Status        int
	Message       string
	OriginalError error
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("HttpError (%d): %s", e.Status, e.Message)
}

func NewHttpError(status int, message string, originalError error) *HttpError {
	return &HttpError{
		Status:        status,
		Message:       message,
		OriginalError: originalError,
	}
}

func NewBadRequestError(message string, originalError error) *HttpError {
	return NewHttpError(http.StatusBadRequest, message, originalError)
}

func NewInternalServerError(messages string, originalError error) *HttpError {
	return NewHttpError(http.StatusInternalServerError, messages, originalError)
}
