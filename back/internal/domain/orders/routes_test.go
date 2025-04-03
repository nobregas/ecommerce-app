package orders

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/middleware/auth"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const testJWTSecret = "secret"

type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) CreateOrderFromCart(userID int, paymentMethod types.PaymentMethod, paymentID string) (*types.OrderHistory, error) {
	args := m.Called(userID, paymentMethod, paymentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.OrderHistory), args.Error(1)
}

func (m *MockOrderService) GetOrdersByUserID(userID int) ([]*types.OrderHistory, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.OrderHistory), args.Error(1)
}

func (m *MockOrderService) GetOrderByID(orderID int) (*types.OrderHistory, error) {
	args := m.Called(orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.OrderHistory), args.Error(1)
}

func (m *MockOrderService) GetOrderWithItems(orderID int) (*types.OrderWithItems, error) {
	args := m.Called(orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.OrderWithItems), args.Error(1)
}

func (m *MockOrderService) GetOrdersWithItems(userID int) ([]*types.OrderWithItems, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.OrderWithItems), args.Error(1)
}

func (m *MockOrderService) UpdateOrderStatus(orderID int, status types.OrderStatus) error {
	args := m.Called(orderID, status)
	return args.Error(0)
}

// MockUserStore é uma implementação mock da interface UserStore
type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) CreateUser(user types.User) error {
	args := m.Called(user)
	return args.Error(0)
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

func createTestToken(userID int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   strconv.Itoa(userID),
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"userRole": string(types.RoleUser),
	})

	signedToken, err := token.SignedString([]byte(testJWTSecret))
	if err != nil {
		panic(err)
	}
	return signedToken
}

func TestCreateOrder(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		payload        string
		mockSetup      func(*MockOrderService, *MockUserStore)
		expectedStatus int
		expectedJSON   bool
	}{
		{
			name:   "Success - Order created",
			userID: 1,
			payload: `{
				"paymentMethod": "CREDIT_CARD",
				"paymentId": "payment123"
			}`,
			mockSetup: func(mos *MockOrderService, mus *MockUserStore) {
				order := &types.OrderHistory{
					ID:            1,
					UserID:        1,
					TotalAmount:   100.0,
					Status:        types.OrderPending,
					PaymentMethod: types.PaymentCreditCard,
					PaymentID:     "payment123",
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				}
				mos.On("CreateOrderFromCart", 1, types.PaymentCreditCard, "payment123").Return(order, nil)
				mus.On("GetUserByID", 1).Return(&types.User{ID: 1}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedJSON:   true,
		},
		{
			name:   "Error - Invalid payment method",
			userID: 1,
			payload: `{
				"paymentMethod": "INVALID_METHOD",
				"paymentId": "payment123"
			}`,
			mockSetup: func(mos *MockOrderService, mus *MockUserStore) {
				mus.On("GetUserByID", 1).Return(&types.User{ID: 1}, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockOrderService)
			mockUserStore := new(MockUserStore)
			tt.mockSetup(mockService, mockUserStore)

			handler := NewHandler(mockService)
			router := mux.NewRouter()
			handler.RegisterRoutes(router, mockUserStore)

			req := httptest.NewRequest("POST", "/orders", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+createTestToken(tt.userID))
			req = req.WithContext(auth.WithUserID(req.Context(), tt.userID))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedJSON {
				var response map[string]interface{}
				err := json.NewDecoder(rr.Body).Decode(&response)
				assert.NoError(t, err)
			}

			mockService.AssertExpectations(t)
			mockUserStore.AssertExpectations(t)
		})
	}
}

func TestGetOrders(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		withItems      bool
		mockSetup      func(*MockOrderService, *MockUserStore)
		expectedStatus int
		expectedJSON   bool
	}{
		{
			name:      "Success - Get orders without items",
			userID:    1,
			withItems: false,
			mockSetup: func(mos *MockOrderService, mus *MockUserStore) {
				orders := []*types.OrderHistory{
					{
						ID:            1,
						UserID:        1,
						TotalAmount:   100.0,
						Status:        types.OrderPending,
						PaymentMethod: types.PaymentCreditCard,
						PaymentID:     "payment123",
						CreatedAt:     time.Now(),
						UpdatedAt:     time.Now(),
					},
				}
				mos.On("GetOrdersByUserID", 1).Return(orders, nil)
				mus.On("GetUserByID", 1).Return(&types.User{ID: 1}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   true,
		},
		{
			name:      "Success - Get orders with items",
			userID:    1,
			withItems: true,
			mockSetup: func(mos *MockOrderService, mus *MockUserStore) {
				ordersWithItems := []*types.OrderWithItems{
					{
						Order: types.OrderHistory{
							ID:            1,
							UserID:        1,
							TotalAmount:   100.0,
							Status:        types.OrderPending,
							PaymentMethod: types.PaymentCreditCard,
							PaymentID:     "payment123",
							CreatedAt:     time.Now(),
							UpdatedAt:     time.Now(),
						},
						Items: []*types.OrderItem{
							{
								OrderID:   1,
								ProductID: 1,
								Quantity:  2,
								Price:     50.0,
							},
						},
					},
				}
				mos.On("GetOrdersWithItems", 1).Return(ordersWithItems, nil)
				mus.On("GetUserByID", 1).Return(&types.User{ID: 1}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockOrderService)
			mockUserStore := new(MockUserStore)
			tt.mockSetup(mockService, mockUserStore)

			handler := NewHandler(mockService)
			router := mux.NewRouter()
			handler.RegisterRoutes(router, mockUserStore)

			url := "/orders"
			if tt.withItems {
				url += "?withItems=true"
			}

			req := httptest.NewRequest("GET", url, nil)
			req.Header.Set("Authorization", "Bearer "+createTestToken(tt.userID))
			req = req.WithContext(auth.WithUserID(req.Context(), tt.userID))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedJSON {
				if tt.withItems {
					var response []*types.OrderWithItems
					err := json.NewDecoder(rr.Body).Decode(&response)
					assert.NoError(t, err)
				} else {
					var response []*types.OrderHistory
					err := json.NewDecoder(rr.Body).Decode(&response)
					assert.NoError(t, err)
				}
			}

			mockService.AssertExpectations(t)
			mockUserStore.AssertExpectations(t)
		})
	}
}

func TestGetOrderByID(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		orderID        string
		withItems      bool
		mockSetup      func(*MockOrderService, *MockUserStore)
		expectedStatus int
		expectedJSON   bool
	}{
		{
			name:      "Success - Get order without items",
			userID:    1,
			orderID:   "1",
			withItems: false,
			mockSetup: func(mos *MockOrderService, mus *MockUserStore) {
				order := &types.OrderHistory{
					ID:            1,
					UserID:        1,
					TotalAmount:   100.0,
					Status:        types.OrderPending,
					PaymentMethod: types.PaymentCreditCard,
					PaymentID:     "payment123",
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				}
				mos.On("GetOrderByID", 1).Return(order, nil)
				mus.On("GetUserByID", 1).Return(&types.User{ID: 1}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   true,
		},
		{
			name:      "Success - Get order with items",
			userID:    1,
			orderID:   "1",
			withItems: true,
			mockSetup: func(mos *MockOrderService, mus *MockUserStore) {
				orderWithItems := &types.OrderWithItems{
					Order: types.OrderHistory{
						ID:            1,
						UserID:        1,
						TotalAmount:   100.0,
						Status:        types.OrderPending,
						PaymentMethod: types.PaymentCreditCard,
						PaymentID:     "payment123",
						CreatedAt:     time.Now(),
						UpdatedAt:     time.Now(),
					},
					Items: []*types.OrderItem{
						{
							OrderID:   1,
							ProductID: 1,
							Quantity:  2,
							Price:     50.0,
						},
					},
				}
				mos.On("GetOrderWithItems", 1).Return(orderWithItems, nil)
				mus.On("GetUserByID", 1).Return(&types.User{ID: 1}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   true,
		},
		{
			name:      "Error - Access denied",
			userID:    1,
			orderID:   "2",
			withItems: false,
			mockSetup: func(mos *MockOrderService, mus *MockUserStore) {
				order := &types.OrderHistory{
					ID:            2,
					UserID:        2, // Outro usuário
					TotalAmount:   100.0,
					Status:        types.OrderPending,
					PaymentMethod: types.PaymentCreditCard,
					PaymentID:     "payment123",
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				}
				mos.On("GetOrderByID", 2).Return(order, nil)
				mus.On("GetUserByID", 1).Return(&types.User{ID: 1}, nil)
			},
			expectedStatus: http.StatusForbidden,
			expectedJSON:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockOrderService)
			mockUserStore := new(MockUserStore)
			tt.mockSetup(mockService, mockUserStore)

			handler := NewHandler(mockService)
			router := mux.NewRouter()
			handler.RegisterRoutes(router, mockUserStore)

			url := "/orders/" + tt.orderID
			if tt.withItems {
				url += "?withItems=true"
			}

			req := httptest.NewRequest("GET", url, nil)
			req.Header.Set("Authorization", "Bearer "+createTestToken(tt.userID))
			req = req.WithContext(auth.WithUserID(req.Context(), tt.userID))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedJSON {
				var response map[string]interface{}
				err := json.NewDecoder(rr.Body).Decode(&response)
				assert.NoError(t, err)
			}

			mockService.AssertExpectations(t)
			mockUserStore.AssertExpectations(t)
		})
	}
}

func TestUpdateOrderStatus(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		orderID        string
		payload        string
		mockSetup      func(*MockOrderService, *MockUserStore)
		expectedStatus int
		expectedJSON   bool
	}{
		{
			name:    "Success - Update order status",
			userID:  1,
			orderID: "1",
			payload: `{
				"status": "COMPLETED"
			}`,
			mockSetup: func(mos *MockOrderService, mus *MockUserStore) {
				order := &types.OrderHistory{
					ID:            1,
					UserID:        1,
					TotalAmount:   100.0,
					Status:        types.OrderPending,
					PaymentMethod: types.PaymentCreditCard,
					PaymentID:     "payment123",
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				}
				mos.On("GetOrderByID", 1).Return(order, nil)
				mos.On("UpdateOrderStatus", 1, types.OrderCompleted).Return(nil)
				mus.On("GetUserByID", 1).Return(&types.User{ID: 1}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   true,
		},
		{
			name:    "Error - Invalid status",
			userID:  1,
			orderID: "1",
			payload: `{
				"status": "INVALID_STATUS"
			}`,
			mockSetup: func(mos *MockOrderService, mus *MockUserStore) {
				order := &types.OrderHistory{
					ID:            1,
					UserID:        1,
					TotalAmount:   100.0,
					Status:        types.OrderPending,
					PaymentMethod: types.PaymentCreditCard,
					PaymentID:     "payment123",
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				}
				mos.On("GetOrderByID", 1).Return(order, nil)
				mus.On("GetUserByID", 1).Return(&types.User{ID: 1}, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   true,
		},
		{
			name:    "Error - Access denied",
			userID:  1,
			orderID: "2",
			payload: `{
				"status": "COMPLETED"
			}`,
			mockSetup: func(mos *MockOrderService, mus *MockUserStore) {
				order := &types.OrderHistory{
					ID:            2,
					UserID:        2, // Outro usuário
					TotalAmount:   100.0,
					Status:        types.OrderPending,
					PaymentMethod: types.PaymentCreditCard,
					PaymentID:     "payment123",
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				}
				mos.On("GetOrderByID", 2).Return(order, nil)
				mus.On("GetUserByID", 1).Return(&types.User{ID: 1}, nil)
			},
			expectedStatus: http.StatusForbidden,
			expectedJSON:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockOrderService)
			mockUserStore := new(MockUserStore)
			tt.mockSetup(mockService, mockUserStore)

			handler := NewHandler(mockService)
			router := mux.NewRouter()
			handler.RegisterRoutes(router, mockUserStore)

			url := "/orders/" + tt.orderID + "/status"

			req := httptest.NewRequest("PATCH", url, strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+createTestToken(tt.userID))
			req = req.WithContext(auth.WithUserID(req.Context(), tt.userID))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedJSON {
				var response map[string]interface{}
				err := json.NewDecoder(rr.Body).Decode(&response)
				assert.NoError(t, err)
			}

			mockService.AssertExpectations(t)
			mockUserStore.AssertExpectations(t)
		})
	}
}
