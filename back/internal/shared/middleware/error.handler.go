package middleware

import (
	"fmt"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/apperrors"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/utils"
	"net/http"
)

func ErrorHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				handleError(w, err)
			}
		}()
		next(w, r)
	}
}

func handleError(w http.ResponseWriter, err interface{}) {
	switch e := err.(type) {
	case *apperrors.AppError:
		writeAppError(w, e)

	case error:
		utils.WriteError(w, http.StatusInternalServerError, e)

	default:
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Something Went Wrong"))

	}
}

func writeAppError(w http.ResponseWriter, e *apperrors.AppError) {
	status := http.StatusInternalServerError
	switch e.Type {
	case apperrors.NotFound:
		status = http.StatusNotFound
	case apperrors.BAD:
		status = http.StatusBadRequest
	case apperrors.Unauthorized:
		status = http.StatusUnauthorized
	}

	utils.WriteJson(w, status, map[string]interface{}{
		"error":   e.Code,
		"message": e.Message,
		"details": e.Details,
	})
}
