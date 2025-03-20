package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) error {
	return WriteJson(w, status, map[string]string{"error": err.Error()})
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

func GetParamIdfromPath(r *http.Request, paramID string) (int, error) {
	// get param id
	vars := mux.Vars(r)
	str, ok := vars[paramID]
	if !ok {
		return -1, fmt.Errorf("missing %s", paramID)
	}

	// convert param id to str
	productID, err := strconv.Atoi(str)
	if err != nil {
		return -1, fmt.Errorf("invalid %s", paramID)
	}
	return productID, nil
}
