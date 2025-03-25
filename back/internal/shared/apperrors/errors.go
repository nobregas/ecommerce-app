package apperrors

import "fmt"

type ErrorType string

const (
	NotFound     ErrorType = "NOT_FOUND"
	BAD          ErrorType = "BAD"
	Unauthorized ErrorType = "UNAUTHORIZED"
	Internal     ErrorType = "INTERNAL"
	FORBIDDEN    ErrorType = "FORBIDDEN"
	CONFLICT     ErrorType = "CONFLICT"
)

type AppError struct {
	Type    ErrorType
	Code    string
	Message string
	Details map[string]interface{}
	Cause   error
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Code)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}
