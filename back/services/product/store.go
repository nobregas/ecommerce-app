package product

import (
	"database/sql"

	"github.com/nobregas/ecommerce-mobile-back/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetProducts() ([]*types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}

	products := make([]*types.Product, 0)
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

func (s *Store) GetProductByID(productID int) (*types.Product, error) {
	return nil, nil
}

func (s *Store) CreateProduct(product types.CreateProductPayload) error {
	_, err := s.db.Exec(
		"INSERT INTO products (title, description, basePrice, stockQuantity) VALUES (?, ?, ?, ?)",
		product.Title, product.Description, product.BasePrice, product.StockQuantity)
	if err != nil {
		return err
	}

	return nil
}

func scanRowsIntoProduct(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)

	err := rows.Scan(
		&product.ID,
		&product.Title,
		&product.Description,
		&product.BasePrice,
		&product.StockQuantity,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return product, nil
}
