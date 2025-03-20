package types

import "time"

type ProductStore interface {
	GetProducts() ([]*Product, error)
	CreateProduct(CreateProductPayload) error
	GetProductByID(productID int) (*Product, error)
	CreateProductWithImages(CreateProductWithImagesPayload) (*Product, error)
	UpdateStock(productID int, quantityChange int) error
	GetInventory(productID int) (*Inventory, error)
	GetImagesForProducts(productIDs []int) (map[int][]ProductImage, error)
}

type Product struct {
	ID          int            `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	BasePrice   float64        `json:"basePrice"`
	Inventory   Inventory      `json:"inventory"`
	Images      []ProductImage `json:"images"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

type Inventory struct {
	ProductID     int `json:"productId"`
	StockQuantity int `json:"stockQuantity"`
	Version       int `json:"-"`
}

type ProductImage struct {
	ID        int    `json:"id"`
	ProductID int    `json:"productId"`
	ImageUrl  string `json:"imageUrl"`
	SortOrder int    `json:"sortOrder"`
}

type CreateProductPayload struct {
	Title         string  `json:"title" validate:"required,min=3,max=100"`
	Description   string  `json:"description" validate:"max=1000"`
	BasePrice     float64 `json:"basePrice" validate:"required,gt=0"`
	StockQuantity int     `json:"stockQuantity" validate:"required,min=0"`
}

type CreateProductWithImagesPayload struct {
	Title         string         `json:"title" validate:"required,min=3,max=100"`
	Description   string         `json:"description" validate:"max=1000"`
	BasePrice     float64        `json:"basePrice" validate:"required,gt=0"`
	StockQuantity int            `json:"stockQuantity" validate:"required,min=0"`
	Images        []ImagePayload `json:"images" validate:"required,min=1,dive"`
}

type ImagePayload struct {
	ImageUrl  string `json:"imageUrl" validate:"required"`
	SortOrder int    `json:"sortOrder" validate:"min=0"`
}
