package apperrors

import "fmt"

func NewEntityNotFound(entity string, id interface{}) *AppError {
	return &AppError{
		Type:    NotFound,
		Code:    "ENTITY_NOT_FOUND",
		Message: fmt.Sprintf("%s not found", entity),
		Details: map[string]interface{}{
			"entity": entity,
			"id":     id,
		},
	}
}

func NewValidationError(field string, reason string) *AppError {
	return &AppError{
		Type:    BAD,
		Code:    "BAD_REQUEST",
		Message: "BAD Request error",
		Details: map[string]interface{}{
			"field":  field,
			"reason": reason,
		},
	}
}

func NewConflictError(field string, reason string) *AppError {
	return &AppError{
		Type:    CONFLICT,
		Code:    "Conflict",
		Message: "Conflict Request error",
		Details: map[string]interface{}{
			"field":  field,
			"reason": reason,
		},
	}
}

// TODO ANOTHER ONES
