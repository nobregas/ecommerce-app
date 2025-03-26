package cart

import (
	"errors"
	"testing"
	"time"

	"github.com/nobregas/ecommerce-mobile-back/internal/shared/apperrors"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCartStore is a mock implementation of the CartStore interface
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

// MockProductStore is a mock implementation of the ProductStore interface
type MockProductStore struct {
	mock.Mock
}

func (m *MockProductStore) GetProductByID(id int) (*types.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Product), args.Error(1)
}

func (m *MockProductStore) CreateProduct(payload types.CreateProductPayload) error {
	args := m.Called(payload)
	return args.Error(0)
}

func (m *MockProductStore) CreateProductWithImages(payload types.CreateProductWithImagesPayload) (*types.Product, error) {
	args := m.Called(payload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Product), args.Error(1)
}

func (m *MockProductStore) DeleteProduct(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductStore) GetImagesForProducts(productIDs []int) (map[int][]types.ProductImage, error) {
	args := m.Called(productIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[int][]types.ProductImage), args.Error(1)
}

func (m *MockProductStore) GetInventory(productID int) (*types.Inventory, error) {
	args := m.Called(productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Inventory), args.Error(1)
}

func (m *MockProductStore) GetProductDetails(productID int, userID int) (*types.ProductDetails, error) {
	args := m.Called(productID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.ProductDetails), args.Error(1)
}

func (m *MockProductStore) GetProducts() ([]*types.Product, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.Product), args.Error(1)
}

func (m *MockProductStore) GetProductsByCategory(categoryID int) ([]*types.Product, error) {
	args := m.Called(categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.Product), args.Error(1)
}

func (m *MockProductStore) GetSimpleProductDetails(productID int) (*[]*types.SimpleProductObject, error) {
	args := m.Called(productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]*types.SimpleProductObject), args.Error(1)
}

func (m *MockProductStore) UpdateProduct(id int, payload types.UpdateProductPayload) error {
	args := m.Called(id, payload)
	return args.Error(0)
}

func (m *MockProductStore) UpdateStock(productID int, quantity int) error {
	args := m.Called(productID, quantity)
	return args.Error(0)
}

// MockProductDiscountStore is a mock implementation of the ProductDiscountStore interface
type MockProductDiscountStore struct {
	mock.Mock
}

func (m *MockProductDiscountStore) GetActiveDiscounts(productID int) ([]*types.ProductDiscount, error) {
	args := m.Called(productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.ProductDiscount), args.Error(1)
}

func (m *MockProductDiscountStore) CreateDiscount(payload *types.CreateProductDiscountPayload) (*types.ProductDiscount, error) {
	args := m.Called(payload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.ProductDiscount), args.Error(1)
}

func (m *MockProductDiscountStore) DeleteDiscount(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductDiscountStore) GetDiscountsByDateRange(productID int, startDate, endDate time.Time) ([]*types.ProductDiscount, error) {
	args := m.Called(productID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.ProductDiscount), args.Error(1)
}

func (m *MockProductDiscountStore) GetDiscountsByProduct(productID int) ([]*types.ProductDiscount, error) {
	args := m.Called(productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.ProductDiscount), args.Error(1)
}

func (m *MockProductDiscountStore) GetDiscoutsByID(id int) (*types.ProductDiscount, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.ProductDiscount), args.Error(1)
}

func (m *MockProductDiscountStore) UpdateDiscount(id int, payload *types.UpdateProductDiscountPayload) (*types.ProductDiscount, error) {
	args := m.Called(id, payload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.ProductDiscount), args.Error(1)
}

func TestServiceCreateCart(t *testing.T) {
	mockCartStore := new(MockCartStore)
	mockProductStore := new(MockProductStore)
	mockDiscountStore := new(MockProductDiscountStore)

	service := NewService(mockCartStore, mockProductStore, mockDiscountStore)

	tests := []struct {
		name          string
		userID        int
		mockSetup     func()
		expectedError error
	}{
		{
			name:   "Success creating cart",
			userID: 1,
			mockSetup: func() {
				mockCartStore.On("CreateCart", 1).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "Error creating cart",
			userID: 1,
			mockSetup: func() {
				mockCartStore.On("CreateCart", 1).Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCartStore.ExpectedCalls = nil
			tt.mockSetup()
			err := service.CreateCart(tt.userID)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			mockCartStore.AssertExpectations(t)
		})
	}
}

func TestServiceGetMyCartItems(t *testing.T) {
	mockCartStore := new(MockCartStore)
	mockProductStore := new(MockProductStore)
	mockDiscountStore := new(MockProductDiscountStore)

	service := NewService(mockCartStore, mockProductStore, mockDiscountStore)

	tests := []struct {
		name          string
		userID        int
		mockSetup     func()
		expectedItems *[]*types.CartItem
		expectedError error
	}{
		{
			name:   "Success getting cart items",
			userID: 1,
			mockSetup: func() {
				items := &[]*types.CartItem{
					{
						CartID:        1,
						ProductID:     1,
						Quantity:      2,
						PriceAtAdding: 10.0,
						AddedAt:       time.Now(),
					},
				}
				mockCartStore.On("GetMyCartItems", 1).Return(items, nil)
			},
			expectedError: nil,
		},
		{
			name:   "Cart not found",
			userID: 1,
			mockSetup: func() {
				mockCartStore.On("GetMyCartItems", 1).Return(nil, apperrors.NewEntityNotFound("cart", 1))
			},
			expectedError: apperrors.NewEntityNotFound("cart", 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCartStore.ExpectedCalls = nil
			tt.mockSetup()
			items, err := service.GetMyCartItems(tt.userID)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, items)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, items)
			}
			mockCartStore.AssertExpectations(t)
		})
	}
}

func TestServiceAddItemToCart(t *testing.T) {
	mockCartStore := new(MockCartStore)
	mockProductStore := new(MockProductStore)
	mockDiscountStore := new(MockProductDiscountStore)

	service := NewService(mockCartStore, mockProductStore, mockDiscountStore)

	tests := []struct {
		name          string
		userID        int
		productID     int
		mockSetup     func()
		expectedError error
	}{
		{
			name:      "Success adding item with discount",
			userID:    1,
			productID: 1,
			mockSetup: func() {
				product := &types.Product{
					ID:        1,
					BasePrice: 100.0,
					Inventory: types.Inventory{StockQuantity: 10},
				}
				discounts := []*types.ProductDiscount{
					{DiscountPercent: 20},
				}
				cartItem := &types.CartItem{
					CartID:        1,
					ProductID:     1,
					Quantity:      1,
					PriceAtAdding: 80.0,
					AddedAt:       time.Now(),
				}

				mockProductStore.On("GetProductByID", 1).Return(product, nil)
				mockDiscountStore.On("GetActiveDiscounts", 1).Return(discounts, nil)
				mockCartStore.On("AddItemToCart", 1, 1, 80.0).Return(cartItem, nil)
			},
			expectedError: nil,
		},
		{
			name:      "Product not found",
			userID:    1,
			productID: 1,
			mockSetup: func() {
				mockProductStore.On("GetProductByID", 1).Return(nil, apperrors.NewEntityNotFound("product", 1))
			},
			expectedError: apperrors.NewEntityNotFound("product", 1),
		},
		{
			name:      "Product out of stock",
			userID:    1,
			productID: 1,
			mockSetup: func() {
				product := &types.Product{
					ID:        1,
					BasePrice: 100.0,
					Inventory: types.Inventory{StockQuantity: 0},
				}
				mockProductStore.On("GetProductByID", 1).Return(product, nil)
			},
			expectedError: apperrors.NewValidationError("product", "product out of stock"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCartStore.ExpectedCalls = nil
			mockProductStore.ExpectedCalls = nil
			mockDiscountStore.ExpectedCalls = nil
			tt.mockSetup()
			item, err := service.AddItemToCart(tt.productID, tt.userID)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, item)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, item)
			}
			mockCartStore.AssertExpectations(t)
			mockProductStore.AssertExpectations(t)
			mockDiscountStore.AssertExpectations(t)
		})
	}
}

func TestServiceRemoveItemFromCart(t *testing.T) {
	mockCartStore := new(MockCartStore)
	mockProductStore := new(MockProductStore)
	mockDiscountStore := new(MockProductDiscountStore)

	service := NewService(mockCartStore, mockProductStore, mockDiscountStore)

	tests := []struct {
		name          string
		userID        int
		productID     int
		mockSetup     func()
		expectedError error
	}{
		{
			name:      "Success removing item",
			userID:    1,
			productID: 1,
			mockSetup: func() {
				mockCartStore.On("RemoveItemFromCart", 1, 1).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:      "Error removing item",
			userID:    1,
			productID: 1,
			mockSetup: func() {
				mockCartStore.On("RemoveItemFromCart", 1, 1).Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCartStore.ExpectedCalls = nil
			tt.mockSetup()
			err := service.RemoveItemFromCart(tt.productID, tt.userID)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			mockCartStore.AssertExpectations(t)
		})
	}
}

func TestServiceGetTotal(t *testing.T) {
	mockCartStore := new(MockCartStore)
	mockProductStore := new(MockProductStore)
	mockDiscountStore := new(MockProductDiscountStore)

	service := NewService(mockCartStore, mockProductStore, mockDiscountStore)

	tests := []struct {
		name          string
		userID        int
		mockSetup     func()
		expectedTotal float64
		expectedError error
	}{
		{
			name:   "Success getting total",
			userID: 1,
			mockSetup: func() {
				mockCartStore.On("GetTotal", 1).Return(100.0, nil)
			},
			expectedTotal: 100.0,
			expectedError: nil,
		},
		{
			name:   "Error getting total",
			userID: 1,
			mockSetup: func() {
				mockCartStore.On("GetTotal", 1).Return(0.0, errors.New("database error"))
			},
			expectedTotal: 0.0,
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCartStore.ExpectedCalls = nil
			tt.mockSetup()
			total, err := service.GetTotal(tt.userID)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Equal(t, tt.expectedTotal, total)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTotal, total)
			}
			mockCartStore.AssertExpectations(t)
		})
	}
}
