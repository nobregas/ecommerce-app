package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/nobregas/ecommerce-mobile-back/internal/shared/apperrors"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var Validate = validator.New()

func ParseJson(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJson(w http.ResponseWriter, status int, v any) {
	if status == http.StatusNoContent || status == http.StatusNotModified {
		w.WriteHeader(status)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		panic(err)
	}
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJson(w, status, map[string]string{"error": err.Error()})
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")

	if tokenAuth != "" {
		tokenAuth = strings.TrimPrefix(tokenAuth, "Bearer ")
	}

	if tokenAuth == "" {
		tokenAuth = r.URL.Query().Get("token")
	}

	return tokenAuth
}

func GetParamIdfromPath(r *http.Request, paramID string) int {
	// get param id
	vars := mux.Vars(r)
	str, ok := vars[paramID]
	if !ok {
		panic(apperrors.NewValidationError(paramID, "invalid paramID"))
	}

	// convert param id to str
	id, err := strconv.Atoi(str)
	if err != nil {
		panic(apperrors.NewValidationError(paramID, "invalid paramID value"))
	}
	return id
}

func FormatValidationError(errs validator.ValidationErrors) error {
	var errorList []string

	for _, err := range errs {
		// if json tag is set, use it
		field := strings.ToLower(err.Field())
		if jsonTag := getJsonTag(err); jsonTag != "" {
			field = jsonTag
		}

		// create friendly message
		switch err.Tag() {
		case "required":
			errorList = append(errorList, fmt.Sprintf("The field %s is required", field))
		case "email":
			errorList = append(errorList, fmt.Sprintf("The field %s must be a valid email", field))
		case "min":
			errorList = append(errorList, fmt.Sprintf("The field %s must have at least %s characters", field, err.Param()))
		case "max":
			errorList = append(errorList, fmt.Sprintf("The field %s must have at most %s characters", field, err.Param()))
		case "eqfield":
			errorList = append(errorList, fmt.Sprintf("The field %s must be equal to %s", field, err.Param()))
		case "cpf":
			errorList = append(errorList, fmt.Sprintf("The field %s must be a valid CPF", field))
		default:
			errorList = append(errorList, fmt.Sprintf("Erro no campo %s: %s", field, err.Tag()))
		}
	}

	return fmt.Errorf("%s", strings.Join(errorList, ". "))
}

func getJsonTag(err validator.FieldError) string {
	if field, ok := err.Type().FieldByName(err.Field()); ok {
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			// remove json tag
			return strings.Split(jsonTag, ",")[0]
		}
	}
	return ""
}

func ParseInt(s string) (int, error) {
	return strconv.Atoi(s)
}
