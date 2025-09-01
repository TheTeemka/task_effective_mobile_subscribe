package merrors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorsToHTTP(err error) int {
	switch {
	case errors.As(err, &validationError):
		return http.StatusBadRequest
	case errors.As(err, &notFoundError):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

type ErrorJson struct {
	Error string `json:"error"`
}

func ErrorToResponseString(err error) string {
	switch {
	case errors.As(err, &validationError) ||
		errors.As(err, &notFoundError):
		return err.Error()
	default:
		return "internal server error"
	}
}

type ValidationError struct {
	message string
}

var validationError *ValidationError

func (e *ValidationError) Error() string {
	return e.message
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{message: message}
}

type NotFoundError struct {
	message string
}

var notFoundError *NotFoundError

func (e *NotFoundError) Error() string {
	return e.message
}

func NewNotFoundErr(message string) *NotFoundError {
	return &NotFoundError{message: message}
}

func GinReturnError(c *gin.Context, err error) {
	status := ErrorsToHTTP(err)
	c.JSON(status, ErrorJson{Error: ErrorToResponseString(err)})
	c.Abort()
}
