package errors

import "net/http"

type HttpError struct {
	HttpCode int
	Message  string
}

// Validation related client error
type ValidationError struct {
	HttpError
}

// NewValidationError return new ValidationError
func NewValidationError(msg string) *ValidationError {
	if msg == "" {
		msg = "Invalid request"
	}

	return &ValidationError{
		HttpError{
			HttpCode: 400,
			Message:  msg,
		},
	}
}

// Error func for returning error message of a validation
func (e *ValidationError) Error() string {
	return e.Message
}

// NotFoundError for error not found
type NotFoundError struct {
	HttpError
}

// NewNotFoundError returns NotFoundError
func NewNotFoundError(msg string) *NotFoundError {
	return &NotFoundError{
		HttpError{
			HttpCode: http.StatusNotFound,
			Message:  msg,
		},
	}
}

// Error func returns message for NotFoundError
func (e *NotFoundError) Error() string {
	return e.Message
}
