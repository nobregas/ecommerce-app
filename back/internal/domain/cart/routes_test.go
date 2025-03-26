package cart

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/middleware/auth"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const testJWTSecret = "secret"

// MockCartService é uma implementação mock da interface CartService
type MockCartService struct {
	mock.Mock
}

func (m *MockCartService) CreateCart(userID int) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockCartService) GetMyCartItems(userID int) (*[]*types.CartItem, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]*types.CartItem), args.Error(1)
}

func (m *MockCartService) AddItemToCart(productID int, userID int) (*types.CartItem, error) {
	args := m.Called(productID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.CartItem), args.Error(1)
}

func (m *MockCartService) RemoveItemFromCart(productID int, userID int) error {
	args := m.Called(productID, userID)
	return args.Error(0)
}

func (m *MockCartService) GetTotal(userID int) (float64, error) {
	args := m.Called(userID)
	return args.Get(0).(float64), args.Error(1)
}

// MockUserStore é uma implementação mock da interface UserStore
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

// createTestToken cria um token JWT válido para os testes
func createTestToken(userID int) string {
	token, _ := auth.CreateJWT([]byte("test-jwt-secret-key"), userID, types.RoleUser)
	return token
}

func TestCreateCart(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		mockSetup      func(*MockCartService, *MockUserStore)
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name:   "Success - Cart created",
			userID: 1,
			mockSetup: func(mcs *MockCartService, mus *MockUserStore) {
				mcs.On("CreateCart", 1).Return(nil)
				mus.On("GetUserByID", 1).Return(&types.User{ID: 1}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]string{
				"message": "Cart created successfully",
			},
		},
		{
			name:   "Error - Service error",
			userID: 1,
			mockSetup: func(mcs *MockCartService, mus *MockUserStore) {
				mcs.On("CreateCart", 1).Return(assert.AnError)
				mus.On("GetUserByID", 1).Return(&types.User{ID: 1}, nil)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCartService)
			mockUserStore := new(MockUserStore)
			tt.mockSetup(mockService, mockUserStore)

			handler := NewHandler(mockService)
			router := mux.NewRouter()
			handler.RegisterRoutes(router, mockUserStore)

			req := httptest.NewRequest("POST", "/cart", nil)
			req.Header.Set("Authorization", "Bearer "+createTestToken(tt.userID))
			req = req.WithContext(auth.WithUserID(req.Context(), tt.userID))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedBody != nil {
				var response map[string]string
				err := json.NewDecoder(rr.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}

			mockService.AssertExpectations(t)
			mockUserStore.AssertExpectations(t)
		})
	}
}

func TestGetMyCartItems(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		mockSetup      func(*MockCartService, *MockUserStore)
		expectedStatus int
		expectedItems  *[]*types.CartItem
	}{
		{
			name:   "Success - Items retrieved",
			userID: 1,
			mockSetup: func(mcs *MockCartService, mus *MockUserStore) {
				items := &[]*types.CartItem{
					{CartID: 1, ProductID: 1, Quantity: 2},
				}
				mcs.On("GetMyCartItems", 1).Return(items, nil)
				mus.On("GetUserByID", 1).Return(&types.User{ID: 1}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedItems: &[]*types.CartItem{
				{CartID: 1, ProductID: 1, Quantity: 2},
			},
		},
		{
			name:   "Error - Service error",
			userID: 1,
			mockSetup: func(mcs *MockCartService, mus *MockUserStore) {
				mcs.On("GetMyCartItems", 1).Return(nil, assert.AnError)
				mus.On("GetUserByID", 1).Return(&types.User{ID: 1}, nil)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCartService)
			mockUserStore := new(MockUserStore)
			tt.mockSetup(mockService, mockUserStore)

			handler := NewHandler(mockService)
			router := mux.NewRouter()
			handler.RegisterRoutes(router, mockUserStore)

			req := httptest.NewRequest("GET", "/cart/items", nil)
			req.Header.Set("Authorization", "Bearer "+createTestToken(tt.userID))
			req = req.WithContext(auth.WithUserID(req.Context(), tt.userID))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedItems != nil {
				var response []*types.CartItem
				err := json.NewDecoder(rr.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, *tt.expectedItems, response)
			}

			mockService.AssertExpectations(t)
			mockUserStore.AssertExpectations(t)
		})
	}
}

func TestAddItemToCart(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		productID      string
		mockSetup      func(*MockCartService, *MockUserStore)
		expectedStatus int
		expectedItem   *types.CartItem
	}{
		{
			name:      "Success - Item added",
			userID:    1,
			productID: "1",
			mockSetup: func(mcs *MockCartService, mus *MockUserStore) {
				item := &types.CartItem{CartID: 1, ProductID: 1, Quantity: 1}
				mcs.On("AddItemToCart", 1, 1).Return(item, nil)
				mus.On("GetUserByID", 1).Return(&types.User{ID: 1}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedItem:   &types.CartItem{CartID: 1, ProductID: 1, Quantity: 1},
		},
		{
			name:      "Error - Invalid product ID",
			userID:    1,
			productID: "invalid",
			mockSetup: func(mcs *MockCartService, mus *MockUserStore) {
				mus.On("GetUserByID", 1).Return(&types.User{ID: 1}, nil)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "Error - Service error",
			userID:    1,
			productID: "1",
			mockSetup: func(mcs *MockCartService, mus *MockUserStore) {
				mcs.On("AddItemToCart", 1, 1).Return(nil, assert.AnError)
				mus.On("GetUserByID", 1).Return(&types.User{ID: 1}, nil)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCartService)
			mockUserStore := new(MockUserStore)
			tt.mockSetup(mockService, mockUserStore)

			handler := NewHandler(mockService)
			router := mux.NewRouter()
			handler.RegisterRoutes(router, mockUserStore)

			req := httptest.NewRequest("POST", "/cart/items/"+tt.productID, nil)
			req.Header.Set("Authorization", "Bearer "+createTestToken(tt.userID))
			req = req.WithContext(auth.WithUserID(req.Context(), tt.userID))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedItem != nil {
				var response types.CartItem
				err := json.NewDecoder(rr.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, *tt.expectedItem, response)
			}

			mockService.AssertExpectations(t)
			mockUserStore.AssertExpectations(t)
		})
	}
}

func TestRemoveItemFromCart(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		productID      string
		mockSetup      func(*MockCartService, *MockUserStore)
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name:      "Success - Item removed",
			userID:    1,
			productID: "1",
			mockSetup: func(mcs *MockCartService, mus *MockUserStore) {
				mcs.On("RemoveItemFromCart", 1, 1).Return(nil)
				mus.On("GetUserByID", 1).Return(&types.User{ID: 1}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]string{
				"message": "Item removed from cart successfully",
			},
		},
		{
			name:      "Error - Invalid product ID",
			userID:    1,
			productID: "invalid",
			mockSetup: func(mcs *MockCartService, mus *MockUserStore) {
				mus.On("GetUserByID", 1).Return(&types.User{ID: 1}, nil)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "Error - Service error",
			userID:    1,
			productID: "1",
			mockSetup: func(mcs *MockCartService, mus *MockUserStore) {
				mcs.On("RemoveItemFromCart", 1, 1).Return(assert.AnError)
				mus.On("GetUserByID", 1).Return(&types.User{ID: 1}, nil)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCartService)
			mockUserStore := new(MockUserStore)
			tt.mockSetup(mockService, mockUserStore)

			handler := NewHandler(mockService)
			router := mux.NewRouter()
			handler.RegisterRoutes(router, mockUserStore)

			req := httptest.NewRequest("DELETE", "/cart/items/"+tt.productID, nil)
			req.Header.Set("Authorization", "Bearer "+createTestToken(tt.userID))
			req = req.WithContext(auth.WithUserID(req.Context(), tt.userID))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedBody != nil {
				var response map[string]string
				err := json.NewDecoder(rr.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}

			mockService.AssertExpectations(t)
			mockUserStore.AssertExpectations(t)
		})
	}
}

func TestGetTotal(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		mockSetup      func(*MockCartService, *MockUserStore)
		expectedStatus int
		expectedTotal  float64
	}{
		{
			name:   "Success - Total retrieved",
			userID: 1,
			mockSetup: func(mcs *MockCartService, mus *MockUserStore) {
				mcs.On("GetTotal", 1).Return(100.0, nil)
				mus.On("GetUserByID", 1).Return(&types.User{ID: 1}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedTotal:  100.0,
		},
		{
			name:   "Error - Service error",
			userID: 1,
			mockSetup: func(mcs *MockCartService, mus *MockUserStore) {
				mcs.On("GetTotal", 1).Return(0.0, assert.AnError)
				mus.On("GetUserByID", 1).Return(&types.User{ID: 1}, nil)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCartService)
			mockUserStore := new(MockUserStore)
			tt.mockSetup(mockService, mockUserStore)

			handler := NewHandler(mockService)
			router := mux.NewRouter()
			handler.RegisterRoutes(router, mockUserStore)

			req := httptest.NewRequest("GET", "/cart/total", nil)
			req.Header.Set("Authorization", "Bearer "+createTestToken(tt.userID))
			req = req.WithContext(auth.WithUserID(req.Context(), tt.userID))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]float64
				err := json.NewDecoder(rr.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTotal, response["total"])
			}

			mockService.AssertExpectations(t)
			mockUserStore.AssertExpectations(t)
		})
	}
}
