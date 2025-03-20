package product

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nobregas/ecommerce-mobile-back/types"
)

func TestProductServiceHandlers(t *testing.T) {
	productStore := &mockProductStore{}
	userStore := &mockUserStore{}
	handler := NewHandler(productStore, userStore)

	t.Run("should handle get products", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/products", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/products", handler.handleGetProducts).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should fail if the product ID is not a number", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/products/abc", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/products/{productID}", handler.handleGetProducts).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should handle get product by ID", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/products/42", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/products/{productID}", handler.handleGetProducts).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should fail creating a product if the payload is missing", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/products", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/products", handler.handleCreateProduct).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should handle creating a product", func(t *testing.T) {
		payload := types.CreateProductPayload{
			Title:         "test",
			BasePrice:     100,
			Description:   "test description",
			StockQuantity: 10,
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/products", handler.handleCreateProduct).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

/*
func TestGetInventory (t *testing.T) {
	db := setupTestDB()
	store := NewStore(db)

	t.Run("should return valid inventory", func(t *testing.T) {
        // Setup
        _, err := db.Exec("INSERT INTO inventory (product_id, stock_quantity) VALUES (1, 100)")
        require.NoError(t, err)

        // Test
        inv, err := store.GetInventory(1)

        // Verify
        assert.NoError(t, err)
        assert.Equal(t, 1, inv.ProductID)
        assert.Equal(t, 100, inv.StockQuantity)
    })

    t.Run("should return error for non-existent product", func(t *testing.T) {
        _, err := store.GetInventory(999)
        assert.ErrorContains(t, err, "not found")
    })
}
*/

type mockProductStore struct{}

func (m *mockProductStore) GetProductByID(productID int) (*types.Product, error) {
	return &types.Product{}, nil
}

func (m *mockProductStore) GetProducts() ([]*types.Product, error) {
	return []*types.Product{}, nil
}

func (m *mockProductStore) CreateProduct(product types.CreateProductPayload) error {
	return nil
}

func (m *mockProductStore) CreateProductWithImages(product types.CreateProductWithImagesPayload) (*types.Product, error) {
	return nil, nil
}

func (m *mockProductStore) GetImagesForProducts(productIDs []int) (map[int][]types.ProductImage, error)

func (m *mockProductStore) UpdateStock(productID int, quantityChange int) error {
	return nil
}
func (m *mockProductStore) GetInventory(productID int) (*types.Inventory, error) {
	return nil, nil
}

type mockUserStore struct{}

func (s *mockUserStore) GetUserByCPF(cpf string) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) GetUserByID(userID int) (*types.User, error) {
	return &types.User{}, nil
}

func (m *mockUserStore) CreateUser(user types.User) error {
	return nil
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return &types.User{}, nil
}
