package types

import "time"

type ProductStore interface {
	GetProducts() ([]*Product, error)
	CreateProduct(CreateProductPayload) error
	GetProductByID(productID int) (*Product, error)
}

type Product struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	BasePrice     float64   `json:"basePrice"`
	StockQuantity int       `json:"stockQuantity"`
	CreatedAt     time.Time `json: "createdAt"`
	UpdatedAt     time.Time `json: "updatedAt"`
}

type CreateProductPayload struct {
	Title         string  `json:"title" validate:"required,min=3,max=100"`
	Description   string  `json:"description" validate:"max=1000"`
	BasePrice     float64 `json:"basePrice" validate:"required,gt=0"`
	StockQuantity int     `json:"stockQuantity" validate:"required,min=0"`
}
