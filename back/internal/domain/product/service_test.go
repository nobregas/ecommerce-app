package product

import (
	"fmt"
	"testing"
	"time"

	"github.com/nobregas/ecommerce-mobile-back/internal/shared/apperrors"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helper function to compare entity not found errors
func assertPanicsWithEntityNotFound(t *testing.T, entity string, id interface{}, fn func()) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Function did not panic as expected")
		}

		// Check if the error type matches the expected one
		appErr, ok := r.(*apperrors.AppError)
		if !ok {
			t.Fatalf("Incorrect panic type: expected *apperrors.AppError, got %T", r)
		}

		// Verify the error fields
		assert.Equal(t, apperrors.NotFound, appErr.Type)
		assert.Equal(t, "ENTITY_NOT_FOUND", appErr.Code)
		assert.Equal(t, fmt.Sprintf("%s not found", entity), appErr.Message)
		assert.Equal(t, entity, appErr.Details["entity"])
		assert.Equal(t, id, appErr.Details["id"])
	}()

	fn()
}

// Helper function to verify panic with generic error
func assertPanicsWithError(t *testing.T, expectedErr error, fn func()) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Function did not panic as expected")
		}

		// Compare error messages
		err, ok := r.(error)
		if !ok {
			t.Fatalf("Incorrect panic type: expected error, got %T", r)
		}

		assert.Equal(t, expectedErr.Error(), err.Error())
	}()

	fn()
}

// Mock for ProductStore
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

func (m *MockProductStore) GetProductByID(productID int) (*types.Product, error) {
	args := m.Called(productID)
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

// Mock para UserStore
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

// Mock para DiscountStore
type MockDiscountStore struct {
	mock.Mock
}

// Implementação dos métodos ausentes para a interface ProductDiscountStore
func (m *MockDiscountStore) CreateDiscount(payload *types.CreateProductDiscountPayload) (*types.ProductDiscount, error) {
	args := m.Called(payload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.ProductDiscount), args.Error(1)
}

func (m *MockDiscountStore) GetActiveDiscountForProduct(productID int) (*types.ProductDiscount, error) {
	args := m.Called(productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.ProductDiscount), args.Error(1)
}

func (m *MockDiscountStore) DeleteDiscount(discountID int) error {
	args := m.Called(discountID)
	return args.Error(0)
}

func (m *MockDiscountStore) GetActiveDiscounts(userID int) ([]*types.ProductDiscount, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.ProductDiscount), args.Error(1)
}

func (m *MockDiscountStore) GetDiscountsByDateRange(userID int, startDate, endDate time.Time) ([]*types.ProductDiscount, error) {
	args := m.Called(userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.ProductDiscount), args.Error(1)
}

func (m *MockDiscountStore) GetDiscountsByProduct(productID int) ([]*types.ProductDiscount, error) {
	args := m.Called(productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.ProductDiscount), args.Error(1)
}

func (m *MockDiscountStore) GetDiscoutsByID(discountID int) (*types.ProductDiscount, error) {
	args := m.Called(discountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.ProductDiscount), args.Error(1)
}

func (m *MockDiscountStore) UpdateDiscount(discountID int, payload *types.UpdateProductDiscountPayload) (*types.ProductDiscount, error) {
	args := m.Called(discountID, payload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.ProductDiscount), args.Error(1)
}

// Mock para RatingStore
type MockRatingStore struct {
	mock.Mock
}

// Implementação dos métodos ausentes para a interface ProductRatingStore
func (m *MockRatingStore) CreateRating(payload *types.CreateProductRatingPayload, userID int, productID int) (*types.ProductRating, error) {
	args := m.Called(payload, userID, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.ProductRating), args.Error(1)
}

func (m *MockRatingStore) GetAverageRatingForProduct(productID int) (float64, error) {
	args := m.Called(productID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockRatingStore) GetRatingByUserAndProduct(userID int, productID int) (*types.ProductRating, error) {
	args := m.Called(userID, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.ProductRating), args.Error(1)
}

func (m *MockRatingStore) DeleteRating(ratingID int) error {
	args := m.Called(ratingID)
	return args.Error(0)
}

func (m *MockRatingStore) GetAverageRating(productID int) (float64, error) {
	args := m.Called(productID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockRatingStore) GetRating(ratingID int) (*types.ProductRating, error) {
	args := m.Called(ratingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.ProductRating), args.Error(1)
}

func (m *MockRatingStore) GetRatingsByProduct(productID int) ([]*types.ProductRating, error) {
	args := m.Called(productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.ProductRating), args.Error(1)
}

func (m *MockRatingStore) GetRatingsByUser(userID int) ([]*types.ProductRating, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.ProductRating), args.Error(1)
}

func (m *MockRatingStore) UpdateRating(ratingID int, payload *types.UpdateProductRatingPayload) (*types.ProductRating, error) {
	args := m.Called(ratingID, payload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.ProductRating), args.Error(1)
}

func TestGetProductByID(t *testing.T) {
	mockProductStore := new(MockProductStore)
	mockUserStore := new(MockUserStore)
	mockDiscountStore := new(MockDiscountStore)
	mockRatingStore := new(MockRatingStore)

	service := NewProductService(mockProductStore, mockUserStore, mockDiscountStore, mockRatingStore)

	t.Run("Success - Returns product", func(t *testing.T) {
		// Set up mock product
		productID := 1
		createdAt := time.Now()
		updatedAt := time.Now()
		mockProduct := &types.Product{
			ID:          productID,
			Title:       "Test Product",
			Description: "Test Description",
			BasePrice:   99.99,
			Inventory: types.Inventory{
				ProductID:     productID,
				StockQuantity: 10,
			},
			Images:     []types.ProductImage{},
			Categories: []types.Category{},
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
		}

		// Configure mock behavior
		mockProductStore.On("GetProductByID", productID).Return(mockProduct, nil)

		// Call the service method
		result := service.GetProductByID(productID)

		// Assert expectations
		assert.NotNil(t, result)
		assert.Equal(t, productID, result.ID)
		assert.Equal(t, "Test Product", result.Title)
		assert.Equal(t, 99.99, result.BasePrice)
		mockProductStore.AssertExpectations(t)
	})

	t.Run("Failure - Product not found", func(t *testing.T) {
		// Set up non-existent product ID
		productID := 999

		// Configure mock behavior
		mockProductStore.On("GetProductByID", productID).Return(nil, fmt.Errorf("product not found"))

		// Test that the function panics with the appropriate error
		assertPanicsWithEntityNotFound(t, "Product", productID, func() {
			service.GetProductByID(productID)
		})

		mockProductStore.AssertExpectations(t)
	})
}

func TestGetProducts(t *testing.T) {
	mockProductStore := new(MockProductStore)
	mockUserStore := new(MockUserStore)
	mockDiscountStore := new(MockDiscountStore)
	mockRatingStore := new(MockRatingStore)

	service := NewProductService(mockProductStore, mockUserStore, mockDiscountStore, mockRatingStore)

	t.Run("Success - Returns products list", func(t *testing.T) {
		// Setup mock products
		mockProducts := []*types.Product{
			{
				ID:          1,
				Title:       "Product 1",
				Description: "Description 1",
				BasePrice:   99.99,
			},
			{
				ID:          2,
				Title:       "Product 2",
				Description: "Description 2",
				BasePrice:   149.99,
			},
		}

		// Configure mock behavior
		mockProductStore.On("GetProducts").Return(mockProducts, nil)

		// Call the service method
		result := service.GetProducts()

		// Assert expectations
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, "Product 1", result[0].Title)
		assert.Equal(t, "Product 2", result[1].Title)
		mockProductStore.AssertExpectations(t)
	})

	t.Run("Failure - Error retrieving products", func(t *testing.T) {
		// Create a separate mock store just for this test
		failMockProductStore := new(MockProductStore)
		failService := NewProductService(failMockProductStore, mockUserStore, mockDiscountStore, mockRatingStore)

		// Configure mock behavior
		mockError := fmt.Errorf("database error")
		failMockProductStore.On("GetProducts").Return(nil, mockError)

		// Define a function that will be executed to capture the panic
		panicFunc := func() {
			failService.GetProducts()
		}

		// Test that the function panics
		assert.Panics(t, panicFunc)

		failMockProductStore.AssertExpectations(t)
	})
}

func TestGetProductDetails(t *testing.T) {
	mockProductStore := new(MockProductStore)
	mockUserStore := new(MockUserStore)
	mockDiscountStore := new(MockDiscountStore)
	mockRatingStore := new(MockRatingStore)

	service := NewProductService(mockProductStore, mockUserStore, mockDiscountStore, mockRatingStore)

	t.Run("Success - Returns product details", func(t *testing.T) {
		// Setup mock user and product
		userID := 1
		productID := 1
		mockUser := &types.User{
			ID:        userID,
			FullName:  "Test User",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
		}

		productDetails := &types.ProductDetails{
			ID:                 &productID,
			Title:              "Test Product",
			Price:              89.99,
			BasePrice:          99.99,
			DiscountPercentage: 10.0,
			Description:        "Test Description",
			IsFavorite:         true,
			AverageRating:      4.5,
			Images:             []types.ProductImage{},
		}

		// Configure mock behavior
		mockUserStore.On("GetUserByID", userID).Return(mockUser, nil)
		mockProductStore.On("GetProductDetails", userID, productID).Return(productDetails, nil)

		// Call the service method
		result := service.GetProductDetails(userID, productID)

		// Assert expectations
		assert.NotNil(t, result)
		assert.Equal(t, productID, *result.ID)
		assert.Equal(t, "Test Product", result.Title)
		assert.Equal(t, 89.99, result.Price)
		assert.Equal(t, 99.99, result.BasePrice)
		assert.Equal(t, 10.0, result.DiscountPercentage)
		assert.True(t, result.IsFavorite)
		assert.Equal(t, 4.5, result.AverageRating)
		mockUserStore.AssertExpectations(t)
		mockProductStore.AssertExpectations(t)
	})

	t.Run("Failure - User not found", func(t *testing.T) {
		// Setup invalid user ID
		userID := 999
		productID := 1

		// Configure mock behavior
		mockUserStore.On("GetUserByID", userID).Return(nil, fmt.Errorf("user not found"))

		// Test that the function panics with the appropriate error
		assertPanicsWithEntityNotFound(t, "user", userID, func() {
			service.GetProductDetails(userID, productID)
		})

		mockUserStore.AssertExpectations(t)
	})

	t.Run("Failure - Product not found", func(t *testing.T) {
		// Setup valid user but invalid product
		userID := 1
		productID := 999
		mockUser := &types.User{
			ID:        userID,
			FullName:  "Test User",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
		}

		// Configure mock behavior
		mockUserStore.On("GetUserByID", userID).Return(mockUser, nil)
		mockProductStore.On("GetProductDetails", userID, productID).Return(nil, fmt.Errorf("product not found"))

		// Test that the function panics with the appropriate error
		assertPanicsWithEntityNotFound(t, "product", productID, func() {
			service.GetProductDetails(userID, productID)
		})

		mockUserStore.AssertExpectations(t)
		mockProductStore.AssertExpectations(t)
	})
}

func TestCreateProductWithImages(t *testing.T) {
	mockProductStore := new(MockProductStore)
	mockUserStore := new(MockUserStore)
	mockDiscountStore := new(MockDiscountStore)
	mockRatingStore := new(MockRatingStore)

	service := NewProductService(mockProductStore, mockUserStore, mockDiscountStore, mockRatingStore)

	t.Run("Success - Creates product with images", func(t *testing.T) {
		// Setup payload for product creation
		payload := types.CreateProductWithImagesPayload{
			Title:         "New Product",
			Description:   "New Product Description",
			BasePrice:     149.99,
			StockQuantity: 50,
			Images: []types.ImagePayload{
				{
					ImageUrl:  "http://example.com/image1.jpg",
					SortOrder: 1,
				},
				{
					ImageUrl:  "http://example.com/image2.jpg",
					SortOrder: 2,
				},
			},
			CategoryIDs: []int{1, 2},
		}

		// Setup mock product to be returned after creation
		createdAt := time.Now()
		updatedAt := time.Now()
		mockProduct := &types.Product{
			ID:          1,
			Title:       payload.Title,
			Description: payload.Description,
			BasePrice:   payload.BasePrice,
			Inventory: types.Inventory{
				ProductID:     1,
				StockQuantity: payload.StockQuantity,
			},
			Images: []types.ProductImage{
				{
					ID:        1,
					ProductID: 1,
					ImageUrl:  payload.Images[0].ImageUrl,
					SortOrder: payload.Images[0].SortOrder,
				},
				{
					ID:        2,
					ProductID: 1,
					ImageUrl:  payload.Images[1].ImageUrl,
					SortOrder: payload.Images[1].SortOrder,
				},
			},
			Categories: []types.Category{
				{
					ID:   1,
					Name: "Category 1",
				},
				{
					ID:   2,
					Name: "Category 2",
				},
			},
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		// Configure mock behavior
		mockProductStore.On("CreateProductWithImages", payload).Return(mockProduct, nil)

		// Call the service method
		result := service.CreateProductWithImages(payload)

		// Assert expectations
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.ID)
		assert.Equal(t, payload.Title, result.Title)
		assert.Equal(t, payload.Description, result.Description)
		assert.Equal(t, payload.BasePrice, result.BasePrice)
		assert.Equal(t, payload.StockQuantity, result.Inventory.StockQuantity)
		assert.Len(t, result.Images, 2)
		assert.Len(t, result.Categories, 2)
		mockProductStore.AssertExpectations(t)
	})

	t.Run("Failure - Invalid payload", func(t *testing.T) {
		// Setup an invalid payload (missing required fields)
		invalidPayload := types.CreateProductWithImagesPayload{
			// Missing Title and BasePrice which are required
			Description: "Invalid Product",
		}

		// The function should panic due to validation error
		assertPanicsWithError(t, apperrors.NewValidationError("invalid payload", "validation error"), func() {
			service.CreateProductWithImages(invalidPayload)
		})
	})

	t.Run("Failure - Store error", func(t *testing.T) {
		// Setup valid payload
		payload := types.CreateProductWithImagesPayload{
			Title:         "New Product",
			Description:   "New Product Description",
			BasePrice:     149.99,
			StockQuantity: 50,
			Images: []types.ImagePayload{
				{
					ImageUrl:  "http://example.com/image1.jpg",
					SortOrder: 1,
				},
			},
			CategoryIDs: []int{1},
		}

		// Configure mock behavior to return error
		storeError := fmt.Errorf("database error during product creation")
		mockProductStore.On("CreateProductWithImages", payload).Return(nil, storeError)

		// Test that the function panics with the appropriate error
		assertPanicsWithError(t, storeError, func() {
			service.CreateProductWithImages(payload)
		})

		mockProductStore.AssertExpectations(t)
	})
}

func TestDeleteProduct(t *testing.T) {
	mockProductStore := new(MockProductStore)
	mockUserStore := new(MockUserStore)
	mockDiscountStore := new(MockDiscountStore)
	mockRatingStore := new(MockRatingStore)

	service := NewProductService(mockProductStore, mockUserStore, mockDiscountStore, mockRatingStore)

	t.Run("Success - Deletes product", func(t *testing.T) {
		// Setup mock product
		productID := 1
		mockProduct := &types.Product{
			ID:          productID,
			Title:       "Test Product",
			Description: "Test Description",
			BasePrice:   99.99,
		}

		// Configure mock behavior
		mockProductStore.On("GetProductByID", productID).Return(mockProduct, nil)
		mockProductStore.On("DeleteProduct", productID).Return(nil)

		// Call the service method - shouldn't panic
		assert.NotPanics(t, func() {
			service.DeleteProduct(productID)
		})

		mockProductStore.AssertExpectations(t)
	})

	t.Run("Failure - Product not found", func(t *testing.T) {
		// Setup non-existent product ID
		productID := 999

		// Configure mock behavior
		mockProductStore.On("GetProductByID", productID).Return(nil, fmt.Errorf("product not found"))

		// Test that the function panics with the appropriate error
		assertPanicsWithEntityNotFound(t, "Product", productID, func() {
			service.DeleteProduct(productID)
		})

		mockProductStore.AssertExpectations(t)
	})

	t.Run("Failure - Error deleting product", func(t *testing.T) {
		// Setup product that will have an error during deletion
		productID := 2
		mockProduct := &types.Product{
			ID:          productID,
			Title:       "Test Product",
			Description: "Test Description",
			BasePrice:   99.99,
		}
		deleteError := fmt.Errorf("error deleting product")

		// Configure mock behavior
		mockProductStore.On("GetProductByID", productID).Return(mockProduct, nil)
		mockProductStore.On("DeleteProduct", productID).Return(deleteError)

		// Test that the function panics with the appropriate error
		assertPanicsWithError(t, deleteError, func() {
			service.DeleteProduct(productID)
		})

		mockProductStore.AssertExpectations(t)
	})
}
