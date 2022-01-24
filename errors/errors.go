package errors

type HttpError struct {
	HttpCode int
	Message  string
}

// Validation related client error
type ValidationError struct {
	HttpError
}

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

func (e *ValidationError) Error() string {
	return e.Message
}
