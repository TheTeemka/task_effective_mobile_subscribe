package merrors

import "net/http"

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
