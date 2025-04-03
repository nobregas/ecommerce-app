package orders

import (
	"errors"
	"testing"
	"time"

	"github.com/nobregas/ecommerce-mobile-back/internal/shared/apperrors"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrderStore é uma implementação mock da interface OrderStore
type MockOrderStore struct {
	mock.Mock
}

func (m *MockOrderStore) CreateOrder(userID int, totalAmount float64, paymentMethod types.PaymentMethod, paymentID string) (*types.OrderHistory, error) {
	args := m.Called(userID, totalAmount, paymentMethod, paymentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.OrderHistory), args.Error(1)
}

func (m *MockOrderStore) AddOrderItems(orderID int, items []*types.OrderItem) error {
	args := m.Called(orderID, items)
	return args.Error(0)
}

func (m *MockOrderStore) GetOrdersByUserID(userID int) ([]*types.OrderHistory, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.OrderHistory), args.Error(1)
}

func (m *MockOrderStore) GetOrderByID(orderID int) (*types.OrderHistory, error) {
	args := m.Called(orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.OrderHistory), args.Error(1)
}

func (m *MockOrderStore) GetOrderItems(orderID int) ([]*types.OrderItem, error) {
	args := m.Called(orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.OrderItem), args.Error(1)
}

func (m *MockOrderStore) UpdateOrderStatus(orderID int, status types.OrderStatus) error {
	args := m.Called(orderID, status)
	return args.Error(0)
}

func (m *MockOrderStore) GetOrdersWithItems(userID int) ([]*types.OrderWithItems, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.OrderWithItems), args.Error(1)
}

func (m *MockOrderStore) GetOrderWithItems(orderID int) (*types.OrderWithItems, error) {
	args := m.Called(orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.OrderWithItems), args.Error(1)
}

// MockCartStore é uma implementação mock da interface CartStore
type MockCartStore struct {
	mock.Mock
}

func (m *MockCartStore) CreateCart(userID int) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockCartStore) GetMyCartItems(userID int) (*[]*types.CartItem, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]*types.CartItem), args.Error(1)
}

func (m *MockCartStore) AddItemToCart(productID int, userID int, price float64) (*types.CartItem, error) {
	args := m.Called(productID, userID, price)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.CartItem), args.Error(1)
}

func (m *MockCartStore) RemoveItemFromCart(productID int, userID int) error {
	args := m.Called(productID, userID)
	return args.Error(0)
}

func (m *MockCartStore) GetTotal(userID int) (float64, error) {
	args := m.Called(userID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockCartStore) GetCartID(userID int) (int, error) {
	args := m.Called(userID)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockCartStore) GetCartItem(userID int, productID int) (*types.CartItem, error) {
	args := m.Called(userID, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.CartItem), args.Error(1)
}

func (m *MockCartStore) RemoveOneItemFromCart(userID int, productID int) error {
	args := m.Called(userID, productID)
	return args.Error(0)
}

func (m *MockCartStore) RemoveItemsFromCart(userID int) error {
	args := m.Called(userID)
	return args.Error(0)
}

// MockProductStore é uma implementação mock da interface ProductStore
type MockProductStore struct {
	mock.Mock
}

func (m *MockProductStore) GetProducts() ([]*types.Product, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.Product), args.Error(1)
}

func (m *MockProductStore) CreateProduct(payload types.CreateProductPayload) error {
	args := m.Called(payload)
	return args.Error(0)
}

func (m *MockProductStore) GetProductByID(id int) (*types.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Product), args.Error(1)
}

func (m *MockProductStore) CreateProductWithImages(payload types.CreateProductWithImagesPayload) (*types.Product, error) {
	args := m.Called(payload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Product), args.Error(1)
}

func (m *MockProductStore) UpdateStock(productID int, quantityChange int) error {
	args := m.Called(productID, quantityChange)
	return args.Error(0)
}

func (m *MockProductStore) GetInventory(productID int) (*types.Inventory, error) {
	args := m.Called(productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Inventory), args.Error(1)
}

func (m *MockProductStore) GetImagesForProducts(productIDs []int) (map[int][]types.ProductImage, error) {
	args := m.Called(productIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[int][]types.ProductImage), args.Error(1)
}

func (m *MockProductStore) UpdateProduct(productID int, payload types.UpdateProductPayload) error {
	args := m.Called(productID, payload)
	return args.Error(0)
}

func (m *MockProductStore) DeleteProduct(productID int) error {
	args := m.Called(productID)
	return args.Error(0)
}

func (m *MockProductStore) GetProductsByCategory(categoryID int) ([]*types.Product, error) {
	args := m.Called(categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.Product), args.Error(1)
}

func (m *MockProductStore) GetProductDetails(userID int, productID int) (*types.ProductDetails, error) {
	args := m.Called(userID, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.ProductDetails), args.Error(1)
}

func (m *MockProductStore) GetSimpleProductDetails(userID int) (*[]*types.SimpleProductObject, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]*types.SimpleProductObject), args.Error(1)
}

func TestCreateOrderFromCart(t *testing.T) {
	mockOrderStore := new(MockOrderStore)
	mockCartStore := new(MockCartStore)
	mockProductStore := new(MockProductStore)

	service := NewService(mockOrderStore, mockCartStore, mockProductStore)

	tests := []struct {
		name          string
		userID        int
		paymentMethod types.PaymentMethod
		paymentID     string
		mockSetup     func()
		expectedOrder *types.OrderHistory
		expectedError error
	}{
		{
			name:          "Success - Create order from cart",
			userID:        1,
			paymentMethod: types.PaymentCreditCard,
			paymentID:     "payment123",
			mockSetup: func() {
				cartItems := &[]*types.CartItem{
					{
						CartID:        1,
						ProductID:     1,
						Quantity:      2,
						PriceAtAdding: 10.0,
						AddedAt:       time.Now(),
					},
				}
				mockCartStore.On("GetMyCartItems", 1).Return(cartItems, nil)
				mockCartStore.On("GetTotal", 1).Return(20.0, nil)

				order := &types.OrderHistory{
					ID:            1,
					UserID:        1,
					TotalAmount:   20.0,
					Status:        types.OrderPending,
					PaymentMethod: types.PaymentCreditCard,
					PaymentID:     "payment123",
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				}
				mockOrderStore.On("CreateOrder", 1, 20.0, types.PaymentCreditCard, "payment123").Return(order, nil)

				mockProductStore.On("UpdateStock", 1, -2).Return(nil)

				var orderItems []*types.OrderItem
				orderItems = append(orderItems, &types.OrderItem{
					OrderID:   1,
					ProductID: 1,
					Quantity:  2,
					Price:     10.0,
				})
				mockOrderStore.On("AddOrderItems", 1, mock.Anything).Return(nil)
				mockCartStore.On("RemoveItemsFromCart", 1).Return(nil)
			},
			expectedOrder: &types.OrderHistory{
				ID:            1,
				UserID:        1,
				TotalAmount:   20.0,
				Status:        types.OrderPending,
				PaymentMethod: types.PaymentCreditCard,
				PaymentID:     "payment123",
			},
			expectedError: nil,
		},
		{
			name:          "Error - Empty cart",
			userID:        1,
			paymentMethod: types.PaymentCreditCard,
			paymentID:     "payment123",
			mockSetup: func() {
				emptyCartItems := &[]*types.CartItem{}
				mockCartStore.On("GetMyCartItems", 1).Return(emptyCartItems, nil)
			},
			expectedOrder: nil,
			expectedError: apperrors.NewValidationError("cart", "cart is empty"),
		},
		{
			name:          "Error - Invalid payment method",
			userID:        1,
			paymentMethod: "INVALID_METHOD",
			paymentID:     "payment123",
			mockSetup:     func() {},
			expectedOrder: nil,
			expectedError: apperrors.NewValidationError("paymentMethod", "invalid payment method: INVALID_METHOD"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrderStore.ExpectedCalls = nil
			mockCartStore.ExpectedCalls = nil
			mockProductStore.ExpectedCalls = nil
			tt.mockSetup()

			order, err := service.CreateOrderFromCart(tt.userID, tt.paymentMethod, tt.paymentID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, order)

				if ve, ok := tt.expectedError.(*apperrors.AppError); ok {
					assert.IsType(t, &apperrors.AppError{}, err)
					actualErr, ok := err.(*apperrors.AppError)
					assert.True(t, ok)
					assert.Equal(t, ve.Details["field"], actualErr.Details["field"])
				} else {
					assert.EqualError(t, err, tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, order)
				assert.Equal(t, tt.expectedOrder.ID, order.ID)
				assert.Equal(t, tt.expectedOrder.UserID, order.UserID)
				assert.Equal(t, tt.expectedOrder.TotalAmount, order.TotalAmount)
				assert.Equal(t, tt.expectedOrder.Status, order.Status)
				assert.Equal(t, tt.expectedOrder.PaymentMethod, order.PaymentMethod)
				assert.Equal(t, tt.expectedOrder.PaymentID, order.PaymentID)
			}

			mockOrderStore.AssertExpectations(t)
			mockCartStore.AssertExpectations(t)
			mockProductStore.AssertExpectations(t)
		})
	}
}

func TestGetOrdersByUserID(t *testing.T) {
	mockOrderStore := new(MockOrderStore)
	mockCartStore := new(MockCartStore)
	mockProductStore := new(MockProductStore)

	service := NewService(mockOrderStore, mockCartStore, mockProductStore)

	tests := []struct {
		name           string
		userID         int
		mockSetup      func()
		expectedOrders []*types.OrderHistory
		expectedError  error
	}{
		{
			name:   "Success - Get orders",
			userID: 1,
			mockSetup: func() {
				orders := []*types.OrderHistory{
					{
						ID:            1,
						UserID:        1,
						TotalAmount:   20.0,
						Status:        types.OrderPending,
						PaymentMethod: types.PaymentCreditCard,
						PaymentID:     "payment123",
						CreatedAt:     time.Now(),
						UpdatedAt:     time.Now(),
					},
				}
				mockOrderStore.On("GetOrdersByUserID", 1).Return(orders, nil)
			},
			expectedOrders: []*types.OrderHistory{
				{
					ID:            1,
					UserID:        1,
					TotalAmount:   20.0,
					Status:        types.OrderPending,
					PaymentMethod: types.PaymentCreditCard,
					PaymentID:     "payment123",
				},
			},
			expectedError: nil,
		},
		{
			name:   "Error - Database error",
			userID: 1,
			mockSetup: func() {
				mockOrderStore.On("GetOrdersByUserID", 1).Return(nil, errors.New("database error"))
			},
			expectedOrders: nil,
			expectedError:  errors.New("error getting orders: database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrderStore.ExpectedCalls = nil
			mockCartStore.ExpectedCalls = nil
			mockProductStore.ExpectedCalls = nil
			tt.mockSetup()

			orders, err := service.GetOrdersByUserID(tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, orders)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, orders)
				assert.Equal(t, len(tt.expectedOrders), len(orders))
				for i, order := range orders {
					assert.Equal(t, tt.expectedOrders[i].ID, order.ID)
					assert.Equal(t, tt.expectedOrders[i].UserID, order.UserID)
					assert.Equal(t, tt.expectedOrders[i].TotalAmount, order.TotalAmount)
					assert.Equal(t, tt.expectedOrders[i].Status, order.Status)
					assert.Equal(t, tt.expectedOrders[i].PaymentMethod, order.PaymentMethod)
					assert.Equal(t, tt.expectedOrders[i].PaymentID, order.PaymentID)
				}
			}

			mockOrderStore.AssertExpectations(t)
		})
	}
}
