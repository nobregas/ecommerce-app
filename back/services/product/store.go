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
	// find products
	row := s.db.QueryRow("SELECT * FROM products WHERE id = ?", productID)
	product, err := scanRowIntoProduct(row)
	if err != nil {
		return nil, err
	}

	// find images
	rows, err := s.db.Query("SELECT * FROM product_images WHERE productId = ?", productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var img types.ProductImage
		err := rows.Scan(&img.ID, &img.ProductID, &img.ImageUrl, &img.SortOrder)
		if err != nil {
			return nil, err
		}
		product.Images = append(product.Images, img)
	}

	return product, nil
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

func (s *Store) CreateProductWithImages(payload types.CreateProductWithImagesPayload) (*types.Product, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// create product
	res, err := tx.Exec(
		"INSERT INTO products (title, description, basePrice, stockQuantity) VALUES (?, ?, ?, ?)",
		payload.Title, payload.Description, payload.BasePrice, payload.StockQuantity,
	)
	if err != nil {
		return nil, err
	}

	productID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	// create product images
	for _, img := range payload.Images {
		_, err := tx.Exec(
			"INSERT INTO product_images (productId, imageUrl, sortOrder) VALUES (?, ?, ?)",
			productID, img.ImageUrl, img.SortOrder,
		)
		if err != nil {
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// find created product
	return s.GetProductByID(int(productID))
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

func scanRowIntoProduct(row *sql.Row) (*types.Product, error) {
	product := new(types.Product)
	err := row.Scan(
		&product.ID,
		&product.Title,
		&product.Description,
		&product.BasePrice,
		&product.StockQuantity,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	return product, err
}
