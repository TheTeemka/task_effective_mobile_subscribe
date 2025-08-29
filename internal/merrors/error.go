package merrors

import "net/http"

func ErrorsToHTTP(err error) int {
	switch err.(type) {
	case *ValidationError:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
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
