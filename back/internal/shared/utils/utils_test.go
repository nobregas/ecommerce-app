package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"required,min=18,max=120"`
}

func TestParseJson(t *testing.T) {
	tests := []struct {
		name          string
		requestBody   string
		expectedError bool
		expected      TestStruct
	}{
		{
			name:          "Valid JSON",
			requestBody:   `{"name":"John Doe","email":"john@example.com","age":30}`,
			expectedError: false,
			expected: TestStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   30,
			},
		},
		{
			name:          "Invalid JSON",
			requestBody:   `{"name":"John Doe","email":}`,
			expectedError: true,
		},
		{
			name:          "Empty body",
			requestBody:   "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			var err error

			if tt.requestBody == "" {
				// Create a request with empty body
				req, err = http.NewRequest("POST", "/test", nil)
			} else {
				// Create a request with the specified body
				req, err = http.NewRequest("POST", "/test", strings.NewReader(tt.requestBody))
			}
			assert.NoError(t, err)

			var result TestStruct
			err = ParseJson(req, &result)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestWriteJson(t *testing.T) {
	tests := []struct {
		name           string
		status         int
		body           interface{}
		expectedStatus int
		expectedBody   string
		expectedHeader string
	}{
		{
			name:           "Normal response",
			status:         http.StatusOK,
			body:           map[string]string{"message": "success"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"success"}`,
			expectedHeader: "application/json; charset=utf-8",
		},
		{
			name:           "No content",
			status:         http.StatusNoContent,
			body:           nil,
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
			expectedHeader: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			WriteJson(rr, tt.status, tt.body)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Check body if expected
			if tt.expectedBody != "" {
				// Remove newline added by json.Encoder
				body := rr.Body.String()
				body = body[:len(body)-1]
				assert.Equal(t, tt.expectedBody, body)
			}

			// Check header if expected
			if tt.expectedHeader != "" {
				assert.Equal(t, tt.expectedHeader, rr.Header().Get("Content-Type"))
			}
		})
	}
}

func TestWriteError(t *testing.T) {
	rr := httptest.NewRecorder()
	err := fmt.Errorf("test error")

	WriteError(rr, http.StatusBadRequest, err)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Equal(t, "test error", response["error"])
}

func TestGetTokenFromRequest(t *testing.T) {
	tests := []struct {
		name          string
		setupRequest  func(*http.Request)
		expectedToken string
	}{
		{
			name: "Token in Authorization header",
			setupRequest: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer test-token")
			},
			expectedToken: "test-token",
		},
		{
			name: "Token in query parameter",
			setupRequest: func(r *http.Request) {
				q := r.URL.Query()
				q.Add("token", "test-token")
				r.URL.RawQuery = q.Encode()
			},
			expectedToken: "test-token",
		},
		{
			name: "No token",
			setupRequest: func(r *http.Request) {
				// No setup needed
			},
			expectedToken: "",
		},
		{
			name: "Authorization header without Bearer prefix",
			setupRequest: func(r *http.Request) {
				r.Header.Set("Authorization", "test-token")
			},
			expectedToken: "test-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/test", nil)
			tt.setupRequest(req)

			token := GetTokenFromRequest(req)
			assert.Equal(t, tt.expectedToken, token)
		})
	}
}

func TestGetParamIdfromPath(t *testing.T) {
	tests := []struct {
		name        string
		paramName   string
		paramValue  string
		expectedID  int
		shouldPanic bool
	}{
		{
			name:        "Valid ID",
			paramName:   "id",
			paramValue:  "123",
			expectedID:  123,
			shouldPanic: false,
		},
		{
			name:        "Invalid ID format",
			paramName:   "id",
			paramValue:  "abc",
			shouldPanic: true,
		},
		{
			name:        "Missing ID",
			paramName:   "missing",
			paramValue:  "",
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := mux.NewRouter()
			router.HandleFunc("/test/{id}", func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if r := recover(); r != nil {
						if !tt.shouldPanic {
							t.Errorf("Should not have panicked but did: %v", r)
						}
					} else if tt.shouldPanic {
						t.Error("Should have panicked but didn't")
					}
				}()

				id := GetParamIdfromPath(r, tt.paramName)
				if !tt.shouldPanic {
					assert.Equal(t, tt.expectedID, id)
				}
			})

			req, _ := http.NewRequest("GET", fmt.Sprintf("/test/%s", tt.paramValue), nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
		})
	}
}

// TestValidationStruct para testes de validação
type TestValidationStruct struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Age      int    `json:"age" validate:"required,min=18,max=120"`
}
