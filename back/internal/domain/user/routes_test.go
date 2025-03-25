package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/middleware/auth"
	types "github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Context key for tests
type contextKey string

const userIDKey contextKey = "userId"

// Custom function to get ID from context in tests
type GetUserIDFromContextFunc func(ctx context.Context) int

// MockUserStore is a mock for the UserStore interface
type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) GetUserByEmail(email string) (*types.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.User), args.Error(1)
}

func (m *MockUserStore) GetUserByID(id int) (*types.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.User), args.Error(1)
}

func (m *MockUserStore) GetUserByCPF(cpf string) (*types.User, error) {
	args := m.Called(cpf)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.User), args.Error(1)
}

func (m *MockUserStore) CreateUser(user types.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// Helper function to create a JWT token for tests
func createTestToken(userID int, role types.UserRole) string {
	// Test key
	secret := []byte("test-secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   fmt.Sprintf("%d", userID),
		"exp":      time.Now().Add(time.Hour).Unix(),
		"userRole": string(role),
	})

	tokenString, _ := token.SignedString(secret)
	return tokenString
}

func TestHandleLogin(t *testing.T) {
	tests := []struct {
		name           string
		payload        types.LoginUserPayload
		setupMock      func(*MockUserStore)
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name: "Success - Valid credentials",
			payload: types.LoginUserPayload{
				Email:    "test@email.com",
				Password: "password123",
			},
			setupMock: func(m *MockUserStore) {
				// Hash for "password123"
				hashedPassword, _ := auth.HashPassword("password123")

				m.On("GetUserByEmail", "test@email.com").Return(&types.User{
					ID:       1,
					Email:    "test@email.com",
					Password: hashedPassword,
					Role:     types.RoleUser,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"token": ""},
		},
		{
			name: "Failure - Email not found",
			payload: types.LoginUserPayload{
				Email:    "notexists@email.com",
				Password: "password123",
			},
			setupMock: func(m *MockUserStore) {
				m.On("GetUserByEmail", "notexists@email.com").Return(nil, fmt.Errorf("user not found"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Failure - Incorrect password",
			payload: types.LoginUserPayload{
				Email:    "test@email.com",
				Password: "wrong_password",
			},
			setupMock: func(m *MockUserStore) {
				hashedPassword, _ := auth.HashPassword("password123")

				m.On("GetUserByEmail", "test@email.com").Return(&types.User{
					ID:       1,
					Email:    "test@email.com",
					Password: hashedPassword,
					Role:     types.RoleUser,
				}, nil)
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := new(MockUserStore)
			tt.setupMock(mockStore)

			handler := NewHandler(mockStore)

			// Create request
			payloadBytes, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			// Record response
			rr := httptest.NewRecorder()

			// Execute handler
			handler.HandleLogin(rr, req)

			// Verify status
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// For success cases, verify token presence
			if tt.expectedStatus == http.StatusOK {
				var response map[string]string
				json.Unmarshal(rr.Body.Bytes(), &response)
				assert.Contains(t, response, "token")
				assert.NotEmpty(t, response["token"])
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestHandleRegister(t *testing.T) {
	tests := []struct {
		name           string
		payload        types.RegisterUserPayload
		setupMock      func(*MockUserStore)
		expectedStatus int
	}{
		{
			name: "Success - Valid registration",
			payload: types.RegisterUserPayload{
				FullName: "Test User",
				Email:    "new@email.com",
				Cpf:      "12345678901",
				Password: "password123",
			},
			setupMock: func(m *MockUserStore) {
				m.On("GetUserByEmail", "new@email.com").Return(nil, fmt.Errorf("user not found"))
				m.On("GetUserByCPF", "12345678901").Return(nil, fmt.Errorf("user not found"))
				m.On("CreateUser", mock.AnythingOfType("types.User")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Failure - Email already registered",
			payload: types.RegisterUserPayload{
				FullName: "Test User",
				Email:    "existing@email.com",
				Cpf:      "12345678901",
				Password: "password123",
			},
			setupMock: func(m *MockUserStore) {
				m.On("GetUserByEmail", "existing@email.com").Return(&types.User{}, nil)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name: "Failure - CPF already registered",
			payload: types.RegisterUserPayload{
				FullName: "Test User",
				Email:    "new@email.com",
				Cpf:      "12345678901",
				Password: "password123",
			},
			setupMock: func(m *MockUserStore) {
				m.On("GetUserByEmail", "new@email.com").Return(nil, fmt.Errorf("user not found"))
				m.On("GetUserByCPF", "12345678901").Return(&types.User{}, nil)
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := new(MockUserStore)
			tt.setupMock(mockStore)

			handler := NewHandler(mockStore)

			// Create request
			payloadBytes, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			// Record response
			rr := httptest.NewRecorder()

			// Execute handler
			handler.HandleRegister(rr, req)

			// Verify status
			assert.Equal(t, tt.expectedStatus, rr.Code)

			mockStore.AssertExpectations(t)
		})
	}
}

// Testable handler that allows injection of the ID retrieval function
type TestableHandler struct {
	store                types.UserStore
	getUserIDFromContext GetUserIDFromContextFunc
}

// Constructor for the testable handler
func NewTestableHandler(store types.UserStore, getUserIDFunc GetUserIDFromContextFunc) *TestableHandler {
	return &TestableHandler{
		store:                store,
		getUserIDFromContext: getUserIDFunc,
	}
}

// HandleGetCurrentUser is a testable version of the original method
func (h *TestableHandler) HandleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserIDFromContext(r.Context())

	user, err := h.store.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	userDTO := types.UserDTO{
		ID:         user.ID,
		FullName:   user.FullName,
		Email:      user.Email,
		Cpf:        user.Cpf,
		ProfileImg: user.ProfileImg,
		CreatedAt:  user.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userDTO)
}

func TestHandleGetCurrentUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		setupMock      func(*MockUserStore)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "Success - Returns current user",
			userID: 1,
			setupMock: func(m *MockUserStore) {
				createdAt := time.Now()
				m.On("GetUserByID", 1).Return(&types.User{
					ID:         1,
					FullName:   "Test User",
					Email:      "test@email.com",
					Cpf:        "12345678901",
					ProfileImg: "https://example.com/img.png",
					CreatedAt:  createdAt,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response types.UserDTO
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, 1, response.ID)
				assert.Equal(t, "Test User", response.FullName)
				assert.Equal(t, "test@email.com", response.Email)
				assert.Equal(t, "12345678901", response.Cpf)
			},
		},
		{
			name:   "Failure - User not found",
			userID: 999,
			setupMock: func(m *MockUserStore) {
				m.On("GetUserByID", 999).Return(nil, fmt.Errorf("user not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := new(MockUserStore)
			tt.setupMock(mockStore)

			// Mocked function to return the test ID
			getUserIDFunc := func(ctx context.Context) int {
				return tt.userID
			}

			handler := NewTestableHandler(mockStore, getUserIDFunc)

			// Create request
			req, _ := http.NewRequest("GET", "/me", nil)

			// Record response
			rr := httptest.NewRecorder()

			// Execute handler
			handler.HandleGetCurrentUser(rr, req)

			// Verify status
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Additional response checks
			if tt.checkResponse != nil {
				tt.checkResponse(t, rr)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestRegisterRoutes(t *testing.T) {
	mockStore := new(MockUserStore)
	handler := NewHandler(mockStore)

	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	// Verify if expected routes were registered
	routes := []struct {
		path   string
		method string
	}{
		{"/login", "POST"},
		{"/register", "POST"},
		{"/me", "GET"},
	}

	for _, route := range routes {
		var match mux.RouteMatch
		req, _ := http.NewRequest(route.method, route.path, nil)
		matched := router.Match(req, &match)
		assert.True(t, matched, "Route %s %s is not registered", route.method, route.path)
	}
}
