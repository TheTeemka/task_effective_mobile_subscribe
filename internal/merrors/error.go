package merrors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorsToHTTP(err error) int {
	switch err.(type) {
	case *ValidationError:
		return http.StatusBadRequest
	case *NotFoundError:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

type ErrorJson struct {
	Error string `json:"error"`
}

func ErrorToResponseString(err error) string {
	switch err.(type) {
	case *ValidationError, *NotFoundError:
		return err.Error()
	default:
		return "internal server error"
	}
}

type ValidationError struct {
	message string
}

func (e *ValidationError) Error() string {
	return e.message
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{message: message}
}

type NotFoundError struct {
	message string
}

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
